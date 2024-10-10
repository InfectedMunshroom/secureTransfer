package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status: %v", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	fmt.Println("File downloaded successfully!")
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: download <server_url>/download?filename=<file name> <download path>")
		return

	}

	url := os.Args[1]
	filename := os.Args[2]

	err := downloadFile(filename, url)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Download completed!")
	}
}
