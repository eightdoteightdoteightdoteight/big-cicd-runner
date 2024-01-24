// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CiCdRequest struct {
	ID       int    `json:"id"`
	RepoName string `json:"repository"`
	CommitID string `json:"ref"`
}

type CdRequest struct {
	ID        int    `json:"id"`
	imageName string `json:"image"`
	tag       string `json:"tag"`
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

func CiCdHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode de la requête est POST
	if r.Method != http.MethodPost {
		pathToYaml := "test.yml"
		stagesExecution(pathToYaml)
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Lire le body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du corps de la requête", http.StatusInternalServerError)
		return
	}

	var requestData CiCdRequest

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Erreur lors du décodage du JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("ID: %d, imageName: %s, tag: %s\n", requestData.ID, requestData.RepoName, requestData.CommitID)

	// Répondre au client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Requête POST traitée avec succès"))

	pathToYaml := ""
	stagesExecution(pathToYaml)
}

func CdHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode de la requête est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Lire le body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du corps de la requête", http.StatusInternalServerError)
		return
	}

	var requestData CdRequest

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Erreur lors du décodage du JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("ID: %d, imageName: %s, tag: %s\n", requestData.ID, requestData.imageName, requestData.tag)

	// Répondre au client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Requête POST traitée avec succès"))
}
