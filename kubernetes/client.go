/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package kubernetes

import (
	"github.com/knadh/koanf"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClientUsing creates a new Kubernetes client using the provided config
func NewClientUsing(config *koanf.Koanf) (kubernetes.Interface, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{}
	loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)

	kubernetesConfig, err := loader.ClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
