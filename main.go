/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package main

import "dolittle.io/fleet-observer/cmd"

func main() {
	cmd.Execute()
	//logger := zerolog.New(os.Stdout)
	//logger.Info().Msg("Starting observer")

	//client, err := kubernetes.NewClientWithDefaultConfig()
	//if err != nil {
	//	return
	//}

	//factory := informers.NewSharedInformerFactory(client, 1*time.Minute)

	//observer := kubernetes.NewObserver("namespaces", factory.Core().V1().Namespaces().Informer(), logger)

	//stop := make(chan struct{})
	//observer.Start(kubernetes.ObserverHandlerFuncs{
	//	HandleFunc: func(obj any) error {
	//		logger.Info().Interface("obj", obj).Msg("Handling")
	//		return nil
	//	},
	//}, stop)

	//go factory.Start(stop)
	//factory.WaitForCacheSync(stop)

	//<-stop
}
