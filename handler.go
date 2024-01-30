package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type CiCdRequestBody struct {
	Id       string `json:"id"`
	RepoName string `json:"repository"`
	Ref      string `json:"Ref"`
}

type CdRequestBody struct {
	Id    string `json:"id"`
	Image string `json:"image"`
	Tag   string `json:"tag"`
}

func CiCdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture du corps de la requête", http.StatusBadRequest)
			return
		}
		fmt.Printf("Corps de la requête : %s\n", string(body))

		var requestData CiCdRequestBody

		err = json.Unmarshal(body, &requestData)

		if err != nil {
			http.Error(w, "Erreur lors du décodage du JSON", http.StatusBadRequest)
			return
		}

		id := requestData.Id
		repoName := requestData.RepoName
		ref := requestData.Ref

		w.WriteHeader(http.StatusOK)

		repoURL := "https://github.com/" + repoName + ".git"
		projectName := strings.Split(repoName, "/")[1]
		tag := strings.Split(ref, "heads/")[1]

		cd_home := exec.Command("cd", "~")
		cd_home_output, err := cd_home.CombinedOutput()
		if err != nil {
			fmt.Println("Erreur lors de l'exécution de la commande:", string(cd_home_output))
		}

		if exists, err := folderExists(repoName); err != nil {
			fmt.Println("Error:", err)
			return
		} else if exists == false {
			cmd := exec.Command("git", "clone", repoURL)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Erreur lors de l'exécution de la commande:", string(output))
			}
		}
		checkout := fmt.Sprintf("origin/%s", ref)
		cmd := exec.Command("cd", repoName, ";", "git", "checkout", checkout, ";", "git", "pull")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Erreur lors de l'exécution de la commande:", string(output))
		}

		pathToYaml := "/big_ci.yml"
		stagesExecution(pathToYaml, id)

		go func() {
			oldImage, err := updateDeployment("imt-framework-staging", projectName, tag)
			if err != nil {
				fmt.Println("Error:", err)
				http.Error(w, "Erreur lors du déploiement", http.StatusInternalServerError)
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
					fmt.Println("Error:", err)
					http.Error(w, "Erreur lors du rollback", http.StatusInternalServerError)
					return
				}
				sendJobResult(id, "Rollback", fmt.Sprintf("Rollback effectué à la version %s", oldImageSplit[1]), "Success")
			}

			finishPipeline(id)
		}()

	} else {
		fmt.Println("Méthode non autorisée")
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func CdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture du corps de la requête", http.StatusBadRequest)
			return
		}

		var requestData CdRequestBody

		err = json.Unmarshal(body, &requestData)
		if err != nil {
			http.Error(w, "Erreur lors du décodage du JSON", http.StatusBadRequest)
			return
		}

		id := requestData.Id
		image := requestData.Image
		tag := requestData.Tag

		w.WriteHeader(http.StatusOK)

		go func() {
			oldImage, err := updateDeployment("imt-framework-staging", image, tag)
			if err != nil {
				fmt.Println("Error:", err)
				http.Error(w, "Erreur lors du déploiement", http.StatusInternalServerError)
				return
			}
			sendJobResult(id, "Déploiement sur Kubernetes", "Déploiement terminé sur imt-framework-staging", "Success")
			time.Sleep(10 * time.Second)

			if health := isHealthy("imt-framework-staging", image); health {
				sendJobResult(id, "Tests de santé", "Les tests de santé ont été réussis", "Success")
			} else {
				sendJobResult(id, "Tests de santé", "Les tests de santé ont échoué", "Failed")
				oldImageSplit := strings.Split(oldImage, ":")
				if _, err := updateDeployment("imt-framework-staging", oldImageSplit[0], oldImageSplit[1]); err != nil {
					fmt.Println("Error:", err)
					http.Error(w, "Erreur lors du rollback", http.StatusInternalServerError)
					return
				}
				sendJobResult(id, "Rollback", fmt.Sprintf("Rollback effectué à la version %s", oldImageSplit[1]), "Success")
			}

			finishPipeline(id)
		}()
	} else {
		fmt.Println("Méthode non autorisée")
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
