package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func updateDeployment(namespace string, image string, tag string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	result, err := deploymentsClient.Get(context.TODO(), fmt.Sprintf("%s-deployment", image), metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	oldImage := result.Spec.Template.Spec.Containers[0].Image

	result.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("registry.nathanaudvard.fr/%s:%s", image, tag)
	_, err = deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
	return oldImage, err
}
