/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package observing

import (
	"dolittle.io/fleet-observer/entities"
	"dolittle.io/fleet-observer/kubernetes"
	"dolittle.io/fleet-observer/storage"
	"fmt"
	"github.com/rs/zerolog"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	listersAppsV1 "k8s.io/client-go/listers/apps/v1"
	listersCoreV1 "k8s.io/client-go/listers/core/v1"
	"strings"
	"time"
)

type PodsHandler struct {
	configurations storage.Configurations
	deployments    storage.Deployments
	events         storage.Events
	configmaps     listersCoreV1.ConfigMapLister
	secrets        listersCoreV1.SecretLister
	replicasets    listersAppsV1.ReplicaSetLister
	logger         zerolog.Logger
}

func NewPodsHandler(configurations storage.Configurations, deployments storage.Deployments, events storage.Events, configmaps listersCoreV1.ConfigMapLister, secrets listersCoreV1.SecretLister, replicasets listersAppsV1.ReplicaSetLister, logger zerolog.Logger) *PodsHandler {
	return &PodsHandler{
		configurations: configurations,
		deployments:    deployments,
		events:         events,
		configmaps:     configmaps,
		secrets:        secrets,
		replicasets:    replicasets,
		logger:         logger,
	}
}

func (ph *PodsHandler) Handle(obj any, deleted bool) error {
	pod, ok := obj.(*coreV1.Pod)
	if !ok {
		return ReceivedWrongType(obj, "Pod")
	}

	logger := ph.logger.With().Str("namespace", pod.GetNamespace()).Str("name", pod.GetName()).Logger()

	tenantID, applicationID, environmentName, microserviceID, ok := GetMicroserviceIdentifiers(pod.ObjectMeta)
	if !ok {
		logger.Trace().Msg("Skipping pod because it is missing microservice identifiers")
		return nil
	}

	_, _, ok = getRuntimeAndHeadContainer(pod.Spec)
	if !ok {
		logger.Trace().Msg("Skipping pod because it does not have a runtime and head container")
		return nil
	}

	var tenantsConfName, dolittleConfName, filesConfName, envConfName, envSecName string
	var hasTenantsConf, hasDolittleConf, hasFilesConf, hasEnvConf, hasEnvSec bool

	for _, volume := range pod.Spec.Volumes {
		if volume.Name == "tenants-config" && volume.ConfigMap != nil {
			tenantsConfName = volume.ConfigMap.Name
			hasTenantsConf = true
		}
		if volume.Name == "dolittle-config" && volume.ConfigMap != nil {
			dolittleConfName = volume.ConfigMap.Name
			hasDolittleConf = true
		}
		if volume.Name == "config-files" && volume.ConfigMap != nil {
			filesConfName = volume.ConfigMap.Name
			hasFilesConf = true
		}
	}
	for _, container := range pod.Spec.Containers {
		if container.Name == "head" {
			for _, source := range container.EnvFrom {
				if source.ConfigMapRef != nil && strings.HasSuffix(source.ConfigMapRef.Name, "-env-variables") {
					envConfName = source.ConfigMapRef.Name
					hasEnvConf = true
				}
				if source.SecretRef != nil && strings.HasSuffix(source.SecretRef.Name, "-secret-env-variables") {
					envSecName = source.SecretRef.Name
					hasEnvSec = true
				}
			}
		}
	}

	if !hasTenantsConf || !hasDolittleConf || !hasFilesConf || !hasEnvConf || !hasEnvSec {
		logger.Trace().Msg("Skipping pod because it is missing configuration references")
		return nil
	}

	tenantsConfig, err := ph.configmaps.ConfigMaps(pod.GetNamespace()).Get(tenantsConfName)
	if err != nil {
		return err
	}
	dolittleConfig, err := ph.configmaps.ConfigMaps(pod.GetNamespace()).Get(dolittleConfName)
	if err != nil {
		return err
	}
	filesConfig, err := ph.configmaps.ConfigMaps(pod.GetNamespace()).Get(filesConfName)
	if err != nil {
		return err
	}
	envConfig, err := ph.configmaps.ConfigMaps(pod.GetNamespace()).Get(envConfName)
	if err != nil {
		return err
	}
	envSecret, err := ph.secrets.Secrets(pod.GetNamespace()).Get(envSecName)
	if err != nil {
		return err
	}

	runtimeConfigHasher := kubernetes.NewConfigHasher()
	runtimeConfigHasher.WriteConfigMap(tenantsConfig)
	runtimeConfigHasher.WriteConfigMap(dolittleConfig)

	customerConfigHasher := kubernetes.NewConfigHasher()
	customerConfigHasher.WriteConfigMap(filesConfig)
	customerConfigHasher.WriteConfigMap(envConfig)
	customerConfigHasher.WriteSecret(envSecret)

	replicaset, err := GetPodOwner(pod, ph.replicasets)
	if err != nil {
		return err
	}

	runtimeConfig := entities.NewRuntimeConfiguration(
		tenantID,
		applicationID,
		environmentName,
		microserviceID,
		runtimeConfigHasher.GetComputedHash(),
	)
	if err := ph.configurations.SetRuntime(runtimeConfig); err != nil {
		return err
	}
	logger.Debug().Interface("config", runtimeConfig).Msg("Updated runtime configuration")

	customerConfig := entities.NewArtifactConfiguration(
		tenantID,
		applicationID,
		environmentName,
		microserviceID,
		customerConfigHasher.GetComputedHash(),
	)
	if err := ph.configurations.SetArtifact(customerConfig); err != nil {
		return err
	}
	logger.Debug().Interface("config", customerConfig).Msg("Updated customer configuration")

	instanceID := entities.NewDeploymentInstanceUID(
		tenantID,
		applicationID,
		environmentName,
		fmt.Sprintf("%v", replicaset.GetGeneration()),
		string(pod.GetUID()),
	)
	instance := entities.NewDeploymentInstance(
		tenantID,
		applicationID,
		environmentName,
		fmt.Sprintf("%v", replicaset.GetGeneration()),
		string(pod.GetUID()),
		pod.GetCreationTimestamp().UTC(),
		ph.stoppedTime(deleted),
		customerConfig,
		runtimeConfig,
		pod.Spec.NodeName,
	)
	if err := ph.deployments.SetInstance(instance); err != nil {
		return err
	}
	logger.Debug().Interface("instance", instance).Msg("Updated deployment instance")

	return ph.handlePodRestarts(instanceID, pod, logger)
}

func (ph *PodsHandler) handlePodRestarts(id entities.DeploymentInstanceUID, pod *coreV1.Pod, logger zerolog.Logger) error {
	platformRestart := ph.getRestartsEventFor(id, pod, true, func(status coreV1.ContainerStatus) bool {
		return status.Name == "runtime"
	})
	if platformRestart.Properties.Count > 0 {
		if err := ph.updateRestartEvent(platformRestart); err != nil {
			return err
		}
		logger.Debug().Interface("event", platformRestart).Msg("Updated event")
	}

	customerRestart := ph.getRestartsEventFor(id, pod, false, func(status coreV1.ContainerStatus) bool {
		return status.Name != "runtime"
	})
	if customerRestart.Properties.Count > 0 {
		if err := ph.updateRestartEvent(customerRestart); err != nil {
			return err
		}
		logger.Debug().Interface("event", customerRestart).Msg("Updated event")
	}

	return nil
}

func (ph *PodsHandler) getRestartsEventFor(id entities.DeploymentInstanceUID, pod *coreV1.Pod, platform bool, predicate func(status coreV1.ContainerStatus) bool) entities.Event {
	event := entities.NewRestartEvent(string(pod.GetUID()), 0, time.Now().UTC(), time.Time{}.UTC(), platform, id)

	for _, status := range pod.Status.ContainerStatuses {
		if status.RestartCount > 0 && (status.State.Terminated != nil || status.LastTerminationState.Terminated != nil) && predicate(status) {
			event.Properties.Count += int(status.RestartCount)

			if status.LastTerminationState.Terminated != nil {
				terminated := status.LastTerminationState.Terminated.FinishedAt.UTC()
				if terminated.Before(event.Properties.FirstTime) {
					event.Properties.FirstTime = terminated
				}
			} else {
				terminated := status.State.Terminated.FinishedAt.UTC()
				if terminated.Before(event.Properties.FirstTime) {
					event.Properties.FirstTime = terminated
				}
			}

			if status.State.Terminated != nil {
				terminated := status.State.Terminated.FinishedAt.UTC()
				if terminated.After(event.Properties.LastTime) {
					event.Properties.LastTime = terminated
				}
			} else {
				terminated := status.LastTerminationState.Terminated.FinishedAt.UTC()
				if terminated.After(event.Properties.LastTime) {
					event.Properties.LastTime = terminated
				}
			}
		}
	}

	return event
}

func (ph *PodsHandler) updateRestartEvent(event entities.Event) error {
	oldEvent, exists, err := ph.events.Get(event.UID)
	if err != nil {
		return err
	}

	if exists && oldEvent.Properties.Count < event.Properties.Count {
		if oldEvent.Properties.FirstTime.Before(event.Properties.FirstTime) {
			event.Properties.FirstTime = oldEvent.Properties.FirstTime
		}
		if oldEvent.Properties.LastTime.After(event.Properties.LastTime) {
			event.Properties.LastTime = oldEvent.Properties.LastTime
		}
	}

	return ph.events.Set(event)
}

func (ph *PodsHandler) stoppedTime(deleted bool) *time.Time {
	if !deleted {
		return nil
	}

	now := time.Now().UTC()
	return &now
}

func GetPodOwner(pod *coreV1.Pod, replicasets listersAppsV1.ReplicaSetLister) (*appsV1.ReplicaSet, error) {
	for _, owner := range pod.GetOwnerReferences() {
		if owner.Kind == "ReplicaSet" {
			if replicaset, err := replicasets.ReplicaSets(pod.GetNamespace()).Get(owner.Name); err == nil {
				return replicaset, nil
			}
		}
	}
	return nil, PodOwnerNotFound
}
