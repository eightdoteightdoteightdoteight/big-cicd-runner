package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func updateDeployment(namespace string, image string, tag string) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	patch := fmt.Sprintf(`{"spec":{"template":{"spec":{"containers":[{"name":"%s","image":"%s:%s"}]}}}}`, image, image, tag)
	_, err = deploymentsClient.Patch(context.TODO(), fmt.Sprintf("%s-deployment", image), types.StrategicMergePatchType, []byte(patch), metav1.PatchOptions{})

	return err
}
