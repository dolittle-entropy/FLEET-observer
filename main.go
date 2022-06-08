package main

import (
	"context"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {
	log.Println("Starting up observer")

	//config, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
	//if err != nil {
	//	log.Fatalln(err)
	//}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		log.Fatalln(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	namespaces, err := client.CoreV1().Namespaces().List(context.Background(), metaV1.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	for _, namespace := range namespaces.Items {
		log.Println("Namespace", namespace.GetName())
	}
}
