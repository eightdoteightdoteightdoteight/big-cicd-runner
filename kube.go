package main

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
	"time"
)

func updateDeployment(namespace string, image string, tag string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	result, err := deploymentsClient.Get(context.TODO(), fmt.Sprintf("%s-deployment", image), metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	oldImage := strings.Split(result.Spec.Template.Spec.Containers[0].Image, "/")[1]

	result.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("registry.nathanaudvard.fr/%s:%s", image, tag)
	_, err = deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
	return oldImage, err
}

func cd(id string, projectName string, tag string) {
	oldImage, err := updateDeployment("imt-framework-staging", projectName, tag)
	if errorAndFinish(err, id, "Déploiement sur Kubernetes", "Erreur lors du déploiement sur Kubernetes") {
		return
	}
	sendJobResult(id, "Déploiement sur Kubernetes", "Déploiement terminé sur imt-framework-staging", "Success")

	time.Sleep(10 * time.Second)

	if health := isHealthy("imt-framework-staging", projectName); health {
		sendJobResult(id, "Tests de santé", "Les tests de santé ont été réussis", "Success")
	} else {
		sendJobResult(id, "Tests de santé", "Les tests de santé ont échoué", "Failed")
		oldImageSplit := strings.Split(oldImage, ":")
		if _, err := updateDeployment("imt-framework-staging", oldImageSplit[0], oldImageSplit[1]); err != nil {
			if errorAndFinish(err, id, "Rollback", "Erreur lors du rollback") {
				return
			}
		}
		sendJobResult(id, "Rollback", fmt.Sprintf("Rollback effectué à la version %s", oldImageSplit[1]), "Success")
		finishPipeline(id, "Failed")
		return
	}

	finishPipeline(id, "Success")
}
