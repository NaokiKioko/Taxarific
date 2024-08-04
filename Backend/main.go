package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Welcome to the home page!")
}

func handleRequests() {
	http.Handle("/home", http.HandlerFunc(home))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}
