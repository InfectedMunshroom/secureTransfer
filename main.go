package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"secureTransfer/internal/server"
)

func main() {
	http.HandleFunc("/upload", server.UploadFiles)
	// http.HandleFunc("/download", server.DownloadFile)

	fmt.Println("Server is running on ", os.Args[1])
	err := http.ListenAndServe(os.Args[1], nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
