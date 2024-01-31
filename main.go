// main.go
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

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

func execCmd(idPipeline string, stage string, errorMsg string, args ...string) []byte {
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if errorAndFinish(err, idPipeline, stage, errorMsg) {
		fmt.Println("Output:", string(output))
		return nil
	}
	return output
}

func errorAndFinish(err error, idPipeline string, stage string, errorMsg string) bool {
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la commande:", err)
		sendJobResult(idPipeline, stage, errorMsg, "Failed")
		finishPipeline(idPipeline, "Failed")
		return true
	}
	return false
}
