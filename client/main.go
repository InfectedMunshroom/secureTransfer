package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"secureTransfer/client/upload"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Parse the uploaded file
		r.ParseMultipartForm(10 << 20) // Limit upload size to 10 MB
		file, handler, err := r.FormFile("file")
		if err != nil {
			fmt.Fprintf(w, "Error retrieving the file: %v", err)
			return
		}
		defer file.Close()

		// Create a new file in the current directory
		filePath := filepath.Join(".", handler.Filename)
		destFile, err := os.Create(filePath)
		if err != nil {
			fmt.Fprintf(w, "Error creating the file: %v", err)
			return
		}
		defer destFile.Close()

		// Write the content to the destination file
		_, err = destFile.ReadFrom(file)
		if err != nil {
			fmt.Fprintf(w, "Error writing to the file: %v", err)
			return
		}

		// Get the absolute file path
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			fmt.Fprintf(w, "Error getting absolute path: %v", err)
			return
		}
		err = upload.UploadFile("http://10.0.2.15:8080/upload", absPath)
		if err != nil {
			fmt.Fprintf(w, "Error uploading file: %v", err)
		}

		// Print the absolute file path in the response
		fmt.Fprintf(w, "File uploaded successfully! Absolute file path: %s", absPath)
	} else {
		// Render the HTML form
		http.ServeFile(w, r, "upload.html")
	}
}

func main() {
	http.HandleFunc("/upload", uploadFile)

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
