package server

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
)

// Placeholder for the actual decrypt function
func decrypt(encryptedData []byte) ([]byte, error) {
	// Implement your decryption logic here
	return encryptedData, nil // Modify this to return the actual decrypted data
}

// Function to check for files in a folder and decrypt them
func checkAndDecrypt(folder string) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Printf("Error reading folder: %v", err)
		return
	}

	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		// Check if the file starts with "aes" (or any other condition)
		if strings.HasPrefix(file.Name(), "aes") {
			filePath := filepath.Join(folder, file.Name())

			// Read the encrypted file
			encryptedData, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Printf("Error reading file %s: %v", file.Name(), err)
				continue
			}

			// Decrypt the file
			decryptedData, err := decrypt(encryptedData)
			if err != nil {
				log.Printf("Error decrypting file %s: %v", file.Name(), err)
				continue
			}

			// Save the decrypted file with a new name
			decryptedFilePath := filepath.Join(folder, "decrypted_"+file.Name())
			err = ioutil.WriteFile(decryptedFilePath, decryptedData, 0644)
			if err != nil {
				log.Printf("Error saving decrypted file %s: %v", decryptedFilePath, err)
				continue
			}

			log.Printf("File %s decrypted successfully to %s", file.Name(), decryptedFilePath)
		}
	}
}

func main() {
	folder := "./aes" // Folder to check for files

	// Create a ticker to check the folder every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check and decrypt files in the folder
			checkAndDecrypt(folder)
		}
	}
}
