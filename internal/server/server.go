package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"secureTransfer/encryptdecrypt"
	"strings"
)

// Placeholder for the actual DecryptAES function
func DecryptAES(encryptedData []byte) ([]byte, error) {
	// Implement your decryption logic here.
	// For example, use AES-256 decryption.
	// Return the decrypted data.
	return encryptedData, nil // Modify this to return the actual decrypted data.
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Fprintf(w, "Error parsing form: %v", err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "Error retrieving the file: %v", err)
		return
	}
	defer file.Close()

	// Read file content
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "Error reading the file: %v", err)
		return
	}

	// Determine the destination directory based on the file name
	var destDir string
	var fileContent []byte

	if strings.HasPrefix(handler.Filename, "aes") {
		destDir = "./aes"

		// Decrypt the file content using DecryptAES
		fileContent, err = encryptdecrypt.DecryptAES("final", fileBytes)
		if err != nil {
			fmt.Fprintf(w, "Error decrypting the file: %v", err)
			return
		}
	} else if strings.HasSuffix(handler.Filename, ".png") {
		destDir = "./file"
		fileContent = fileBytes // No decryption, just save the original content
	} else {
		fmt.Fprintf(w, "File must either start with 'aes' or end with '.png'")
		return
	}

	// Create the directory if it does not exist
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err = os.Mkdir(destDir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(w, "Error creating directory: %v", err)
			return
		}
	}

	// Save the decrypted file content
	destPath := filepath.Join(destDir, handler.Filename)
	err = ioutil.WriteFile(destPath, fileContent, 0644)
	if err != nil {
		fmt.Fprintf(w, "Error saving the file: %v", err)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %v\n", handler.Filename)
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Filename is missing", http.StatusBadRequest)
		return
	}

	// Determine the correct directory
	var filepath string
	if strings.HasPrefix(filename, "aes") {
		filepath = "./aes/" + filename
	} else if strings.HasSuffix(filename, ".png") {
		filepath = "./file/" + filename
	} else {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	io.Copy(w, file)
}
