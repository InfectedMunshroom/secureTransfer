package main

import (
	"hello/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/encrypt", handlers.EncryptHandler)
	http.HandleFunc("/decrypt", handlers.DecryptHandler)
	http.Handle("/", http.FileServer(http.Dir("./static/")))

	log.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
