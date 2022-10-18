/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package kubernetes

import (
	"github.com/rs/zerolog"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Observer struct {
	queue  workqueue.RateLimitingInterface
	index  cache.Indexer
	logger zerolog.Logger
}

func NewObserver(name string, informer cache.SharedIndexInformer, logger zerolog.Logger) *Observer {
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), name)
	logger = logger.With().Str("observer", name).Logger()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: queue.Add,
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldMeta, oldOk := oldObj.(metaV1.Object)
			newMeta, newOk := newObj.(metaV1.Object)

			if oldOk && newOk && oldMeta.GetResourceVersion() == newMeta.GetResourceVersion() {
				logger.Trace().Str("namespace", newMeta.GetNamespace()).Str("name", newMeta.GetName()).Msg("Skipping because resource version has not changed")
				return
			}

			queue.Add(newObj)
		},
		DeleteFunc: queue.Add,
	})

	index := informer.GetIndexer()

	return &Observer{
		queue:  queue,
		index:  index,
		logger: logger,
	}
}

func (o *Observer) Start(handler ObserverHandler, stopCh <-chan struct{}) {
	o.logger.Info().Msg("Starting observer")
	go o.handleQueue(handler)
	go o.shutdownWhenStopped(stopCh)
}

func (o *Observer) handleQueue(handler ObserverHandler) {
	for {
		item, shutdown := o.queue.Get()
		if shutdown {
			break
		}

		meta, ok := item.(metaV1.Object)
		if !ok {
			o.logger.Warn().Msg("Will skip handling of object without metadata")
		}

		namespace := meta.GetNamespace()
		name := meta.GetName()

		logger := o.logger.With().Str("namespace", namespace).Str("name", name).Logger()
		logger.Debug().Msg("Handling item")

		key := name
		if len(namespace) > 0 {
			key = namespace + "/" + name
		}

		_, exists, err := o.index.GetByKey(key)
		if err != nil {
			o.queue.AddRateLimited(item)
			logger.Warn().Err(err).Msg("Failed to get item from index")
		} else if err := handler.Handle(item, !exists); err != nil {
			o.queue.AddRateLimited(item)
			logger.Warn().Err(err).Msg("Error occurred while handling item")
		} else {
			o.queue.Forget(item)
			logger.Debug().Msg("Done handling item")
		}

		o.queue.Done(item)
	}

	o.logger.Debug().Msg("Queue has been shut down")
}

func (o *Observer) shutdownWhenStopped(stopCh <-chan struct{}) {
	<-stopCh
	o.logger.Debug().Msg("Stopping queue")
	o.queue.ShutDown()
}

type ObserverHandler interface {
	Handle(obj any, deleted bool) error
}

type ObserverHandlerFuncs struct {
	HandleFunc func(obj any, deleted bool) error
}

func (fs ObserverHandlerFuncs) Handle(obj any, deleted bool) error {
	if fs.HandleFunc != nil {
		return fs.HandleFunc(obj, deleted)
	}
	return nil
}
