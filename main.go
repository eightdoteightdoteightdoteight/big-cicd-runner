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

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")

	pipeline, err := readYaml("test.yml")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	print(pipeline.StagesList[0]) // should print stage1
}

// HelloHandler handles requests to the "/hello" endpoint
func HwHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}
