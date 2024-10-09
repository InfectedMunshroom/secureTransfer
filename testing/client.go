package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func uploadFile(filename string, targetUrl string) error {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Create a new buffer and multipart writer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Create a form file field in the multipart request
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("could not create form file: %v", err)
	}

	// Copy the file content into the form file
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("could not close writer: %v", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", targetUrl, body)
	if err != nil {
		return fmt.Errorf("could not create request: %v", err)
	}

	// Set the Content-Type header to multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request to the server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	// Read and print the server's response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %v", err)
	}

	fmt.Println("Response from server:", string(respBody))
	return nil
}

func main() {
	targetUrl := "http://10.0.2.3:8080/upload"
	filename := "yourfile.txt" // Replace with your actual file path
	err := uploadFile(filename, targetUrl)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("File uploaded successfully!")
	}
}

