package kubernetes

import (
	"errors"
	"github.com/rs/zerolog"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Observer struct {
	queue  workqueue.RateLimitingInterface
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
	})

	return &Observer{
		queue:  queue,
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

		logger := o.logger
		if meta, ok := item.(metaV1.Object); ok {
			logger = logger.With().Str("namespace", meta.GetNamespace()).Str("name", meta.GetName()).Logger()
		}

		logger.Debug().Msg("Handling item")

		if err := handler.Handle(item); err == nil {
			o.queue.Forget(item)
			logger.Debug().Msg("Done handling item")
		} else if errors.Is(err, IrrecoverableError) {
			o.queue.Forget(item)
			logger.Error().Err(err).Msg("Giving up because of fatal error")
		} else {
			o.queue.AddRateLimited(item)
			logger.Warn().Err(err).Msg("Error occurred while handling item")
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
	Handle(obj any) error
}

type ObserverHandlerFuncs struct {
	HandleFunc func(obj any) error
}

func (fs ObserverHandlerFuncs) Handle(obj any) error {
	if fs.HandleFunc != nil {
		return fs.HandleFunc(obj)
	}
	return nil
}
