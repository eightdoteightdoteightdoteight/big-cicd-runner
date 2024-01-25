// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/pipelines/cicd", CiCdHandler)
	mux.HandleFunc("/v1/pipelines/cd", CdHandler)

	var port int = 8080

	fmt.Printf("Server is running on port :%d...\n", port)

	errWeb := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if errWeb != nil {
		fmt.Println("Error:", errWeb)
	}
}

func folderExists(folderPath string) (bool, error) {
	_, err := os.Stat(folderPath)

	if err == nil {
		return true, nil // Folder exists
	}

	if os.IsNotExist(err) {
		return false, nil // Folder doesn't exist
	}

	return false, err
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
		var requestBody CdRequestBody
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestBody); err != nil {
			fmt.Printf("Erreur de décodage JSON : %v\n", err)
			http.Error(w, "Erreur lors du décodage du corps JSON", http.StatusBadRequest)
			return
		}
		fmt.Printf("Valeurs après décodage : %+v\n", requestBody)

		id := requestBody.Id
		image := requestBody.Image
		tag := requestBody.Tag

		fmt.Printf("Pipeline id: %s\n", id)
		fmt.Printf("Image: %s\n", image)
		fmt.Printf("Tag: %s\n", tag)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Requête POST traitée avec succès"))
	} else {
		fmt.Println("Méthode non autorisée")
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
