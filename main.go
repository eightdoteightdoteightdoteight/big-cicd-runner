// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type CiCdRequest struct {
	ID       string `json:"id"`
	RepoName string `json:"repository"`
	CommitID string `json:"ref"`
}

type CdRequest struct {
	ID        string `json:"id"`
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
	// Vérifier que la méthode de la requête est un POST
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

	pathToYaml := "test.yml"
	stagesExecution(pathToYaml, requestData.ID)

	if exists, err := folderExists(requestData.RepoName); err != nil {
		fmt.Println("Error:", err)
	} else if exists {
		fmt.Printf("Le dossier %s existe.\n (faut faire un git fetch -all)", requestData.RepoName)
	} else {
		fmt.Printf("Le dossier %s n'existe pas.\n(faut faire un git clone)", requestData.RepoName)
	}
}

func CdHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode de la requête est un POST
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
