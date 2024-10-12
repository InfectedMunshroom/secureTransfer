package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Helper function to create a multipart file field
func createFileField(writer *multipart.Writer, fieldname, filename string) error {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create a form file field in the multipart writer
	part, err := writer.CreateFormFile(fieldname, filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("error creating form file field: %v", err)
	}

	// Copy the file content to the form file field
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("error copying file content: %v", err)
	}

	return nil
}

// UploadFiles uploads both the .png file and the aes_ .png file to the server
func UploadFiles(filePath, aesFilePath string) error {
	// Prepare a buffer to hold the multipart data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add the main file (<file name>.png) to the request
	if err := createFileField(writer, "file", filePath); err != nil {
		return fmt.Errorf("error adding file field: %v", err)
	}

	// Add the AES-encrypted file (aes_<file name>.png) to the request
	if err := createFileField(writer, "aesFile", aesFilePath); err != nil {
		return fmt.Errorf("error adding AES file field: %v", err)
	}

	// Close the multipart writer to finalize the form data
	err := writer.Close()
	if err != nil {
		return fmt.Errorf("error closing multipart writer: %v", err)
	}

	// Send the POST request to the server
	uploadURL := "http://localhost:8080/upload" // Adjust the URL if necessary
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the content type to multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response from the server
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	// Print the server response
	fmt.Println("Server Response:", string(respBody))
	return nil
}

func main() {
	// Paths to the files to be uploaded
	filePath := os.Args[1]    // Path to the .png file
	aesFilePath := os.Args[2] // Path to the aes_ .png file

	// Upload the files
	err := UploadFiles(filePath, aesFilePath)
	if err != nil {
		fmt.Printf("Error uploading files: %v\n", err)
	} else {
		fmt.Println("Files uploaded successfully.")
	}
}
