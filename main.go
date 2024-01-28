// main.go
package main

import (
	"fmt"
	"net/http"
	"os"
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
