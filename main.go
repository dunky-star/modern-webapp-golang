package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	portNumber = ":3000"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/v1/health", healthCheckHandler)
	http.HandleFunc("/v1/about", aboutUsHandler)

	fmt.Printf("Server started on port %s\n", portNumber)
	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	n, err := fmt.Fprintf(w, "Health check")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Number of bytes written: %d\n", n)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the home page")
}

func aboutUsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the about page")
}
