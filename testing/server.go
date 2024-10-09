package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // Max file size of 10 MB
	if err != nil {
		fmt.Fprintf(w, "Error parsing form: %v", err)
		return
	}

	// Get the file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "Error retrieving the file: %v", err)
		return
	}
	defer file.Close()

	// Create the destination directory if it doesn't exist
	destDir := "./file"
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err = os.Mkdir(destDir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(w, "Error creating directory: %v", err)
			return
		}
	}

	// Create the destination file
	destPath := filepath.Join(destDir, handler.Filename)
	dst, err := os.Create(destPath)
	if err != nil {
		fmt.Fprintf(w, "Error creating the file: %v", err)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Fprintf(w, "Error saving the file: %v", err)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %v\n", handler.Filename)
}

func main() {
	http.HandleFunc("/upload", uploadFile)
	fmt.Println("Server started on 10.0.2.15:8080...")
	err := http.ListenAndServe("10.0.2.15:8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}

