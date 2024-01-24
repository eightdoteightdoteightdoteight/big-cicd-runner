// main.go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", RootHandler)
	mux.HandleFunc("/hello-world", HwHandler)

	var port int = 8080

	fmt.Printf("Server is running on port :%d...\n", port)

	var pathToYaml string = "test.yml"
	stagesExecution(pathToYaml)

	errWeb := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if errWeb != nil {
		fmt.Println("Error:", errWeb)
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

// HelloHandler handles requests to the "/hello" endpoint
func HwHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}
