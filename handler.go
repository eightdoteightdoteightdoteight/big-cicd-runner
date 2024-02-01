package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strings"
)

type CiCdRequestBody struct {
	Id       string `json:"id"`
	RepoName string `json:"name"`
	Url      string `json:"url"`
	Ref      string `json:"ref"`
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
		repoURL := requestData.Url + ".git"
		ref := requestData.Ref

		w.WriteHeader(http.StatusOK)

		tag := strings.Split(ref, "heads/")[1]

		usr, err := user.Current()
		if errorAndFinish(err, id, "Setup", "Erreur lors de la récupération du code source") {
			return
		}

		err = os.Chdir(usr.HomeDir)
		if errorAndFinish(err, id, "Setup", "Erreur lors de la récupération du code source") {
			return
		}

		if exists, err := folderExists(repoName); err != nil {
			if errorAndFinish(err, id, "Setup", "Erreur lors de la récupération du code source") {
				return
			}
		} else if exists == false {
			if output := execCmd(id, "Setup", "Erreur lors de la récupération du code source", "git", "clone", repoURL); output == nil {
				return
			}
		}

		err = os.Chdir(repoName)
		if errorAndFinish(err, id, "Setup", "Erreur lors de la récupération du code source") {
			return
		}

		if output := execCmd(id, "Setup", "Erreur lors de la récupération du code source", "git", "fetch"); output == nil {
			return
		}
		checkout := fmt.Sprintf("origin/%s", tag)
		if output := execCmd(id, "Setup", "Erreur lors de la récupération du code source", "git", "checkout", checkout); output == nil {
			return
		}
		if output := execCmd(id, "Setup", "Erreur lors de la récupération du code source", "git", "pull", "origin", tag); output == nil {
			return
		}

		pathToYaml := "big_ci.yml"
		go stagesExecution(pathToYaml, id, repoName, tag)
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

		go cd(id, image, tag)
	} else {
		fmt.Println("Méthode non autorisée")
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
