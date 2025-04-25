package main

import (
	"log"
	"os"
	"strconv"

	"net/http"
)

const (
	defaultPort = "8000"
	numRows     = 12
	numCols     = 12
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("Starting server on port %s...\n", port)
	http.ListenAndServe(":"+port, newRouter())
}
