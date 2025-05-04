package main

import (
	"log"
	"os"
	"strconv"

	"net/http"

	"github.com/joho/godotenv"
)

const (
	defaultPort     = "8000"
	numRows         = 12
	numCols         = 12
	blockSize       = 40
	mobileBlockSize = 30
	gapSize         = 2
)

func main() {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, newRouter()))
}
