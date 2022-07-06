package main

import (
	"dolittle.io/fleet-observer/kubernetes"
	"github.com/rs/zerolog"
	"k8s.io/client-go/informers"
	"os"
	"time"
)

func main() {
	logger := zerolog.New(os.Stdout)
	logger.Info().Msg("Starting observer")

	client, err := kubernetes.NewClientWithDefaultConfig()
	if err != nil {
		return
	}

	factory := informers.NewSharedInformerFactory(client, 1*time.Minute)

	observer := kubernetes.NewObserver("namespaces", factory.Core().V1().Namespaces().Informer(), logger)

	stop := make(chan struct{})
	observer.Start(kubernetes.ObserverHandlerFuncs{
		HandleFunc: func(obj any) error {
			logger.Info().Interface("obj", obj).Msg("Handling")
			return nil
		},
	}, stop)

	go factory.Start(stop)
	factory.WaitForCacheSync(stop)

	<-stop
}
