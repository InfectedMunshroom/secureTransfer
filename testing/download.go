package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadFile(filepath string, url string) error {
	// Create the file to save the downloaded content
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer out.Close()

	// Send the HTTP GET request to download the file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the server returned a success status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status: %v", resp.Status)
	}

	// Copy the content from the response to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	fmt.Println("File downloaded successfully!")
	return nil
}

func main() {
	if len(os.Args)<3{
		fmt.Println("Usage: download <server_url>/download?filename=<file name> <download path>")
		return

	}

	url := os.Args[1]
	filename := os.Args[2]

	err := downloadFile(filename,url)


	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Download completed!")
	}
}

