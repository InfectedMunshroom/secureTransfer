package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"secureTransfer/encryptdecrypt" // Replace with your actual package path
)

// Structure to capture the response from the server
type DownloadResponse struct {
	AESKey   []byte `json:"aes_key"`
	FileData []byte `json:"file_data"`
}

func main() {
	// Specify the URL of the server (adjust the port and URL if necessary)
	url := "http://localhost:8081/download?file=dec_encrypted_output.png"

	// Send the HTTP GET request to download the file
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error during HTTP request:", err)
		return
	}
	defer response.Body.Close()

	// Check for non-200 status code, which indicates an error from the server
	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Server returned status %d: %s\n", response.StatusCode, string(body))
		return
	}

	// Read and parse the JSON response
	var downloadResponse DownloadResponse
	err = json.NewDecoder(response.Body).Decode(&downloadResponse)
	if err != nil {
		// Print out the response for debugging
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Error decoding JSON response:", err)
		fmt.Println("Raw response body:", string(body))
		return
	}

	// Decrypt the AES key using RSA (Assuming DecryptAES is your decryption function)
	decryptedAESKey, err := encryptdecrypt.DecryptAES("./../final", downloadResponse.AESKey)
	if err != nil {
		fmt.Println("Error in decoding the AES key:", err)
		return
	}
	fmt.Println("Decrypted AES key:", string(decryptedAESKey))

	// Write the encrypted file data to a temporary file named "temp"
	tempFilePath := "temp"
	err = ioutil.WriteFile(tempFilePath, downloadResponse.FileData, 0644)
	if err != nil {
		fmt.Println("Error writing encrypted data to 'temp':", err)
		return
	}
	fmt.Println("Encrypted file data written to 'temp'.")

	// Decrypt the file data using the decrypted AES key (Assuming DecodeFile is your decryption function)
	decryptedFileData, err := encryptdecrypt.DecodeFile(decryptedAESKey, tempFilePath)
	if err != nil {
		fmt.Println("Error decrypting file:", err)
		return
	}

	// Prompt the user to specify the download path
	fmt.Print("Enter the path where you want to save the decrypted file (including file name): ")
	var downloadPath string
	fmt.Scan(&downloadPath)

	// Write the decrypted file data to the specified location
	cleanPath := filepath.Clean(downloadPath)
	err = ioutil.WriteFile(cleanPath, decryptedFileData, 0644)
	if err != nil {
		fmt.Println("Error saving the decrypted file:", err)
		return
	}

	fmt.Println("File downloaded and decrypted successfully at:", cleanPath)
}
