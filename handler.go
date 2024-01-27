package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

type CiCdRequestBody struct {
	ID       string `json:"id"`
	RepoName string `json:"repository"`
	CommitID string `json:"ref"`
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

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Requête POST traitée avec succès"))

		repoURL := "https://github.com/" + requestData.RepoName + ".git"

		if exists, err := folderExists(requestData.RepoName); err != nil {
			fmt.Println("Error:", err)
			return
		} else if exists {
			cmd := exec.Command("git", "-C", requestData.RepoName, "fetch", "--all")
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Erreur lors de l'exécution de la commande:", string(output))
			}
		} else {
			cmd := exec.Command("git", "clone", repoURL)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Erreur lors de l'exécution de la commande:", string(output))
			}
		}

		pathToYaml := requestData.RepoName + "/big_ci.yml"
		stagesExecution(pathToYaml, requestData.ID)

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

		if err := updateDeployment("imt-framework-staging", image, tag); err != nil {
			fmt.Println("Error:", err)
			http.Error(w, "Erreur lors du déploiement", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Requête POST traitée avec succès"))

		sendJobResult(id, "Deploy", "Déploiement terminé sur imt-framework-staging", "Success")
		finishPipeline(id)
	} else {
		fmt.Println("Méthode non autorisée")
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
