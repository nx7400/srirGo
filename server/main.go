package main

import (
	"log"
	"net/http"
)

// Main function of server.
func main() {

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
