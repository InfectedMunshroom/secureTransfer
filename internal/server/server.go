package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"secureTransfer/encryptdecrypt"
	"strings"
)

// Placeholder function to decrypt the AES key with the server's private key
func decryptAESKey(encryptedKey []byte) ([]byte, error) {
	// Implement your decryption logic for the AES key here
	return encryptedKey, nil // Replace with actual decrypted AES key
}

// Placeholder function to decrypt the file with the AES key
func decryptWithAESKey(encryptedData, aesKey []byte) ([]byte, error) {
	// Implement your AES decryption logic here
	return encryptedData, nil // Replace with actual decrypted file data
}

func UploadFiles(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // Allow file size up to 10 MB
	if err != nil {
		fmt.Fprintf(w, "Error parsing form: %v", err)
		return
	}

	// Retrieve files from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "Error retrieving file: %v", err)
		return
	}
	defer file.Close()

	aesFile, aesHandler, err := r.FormFile("aesFile")
	if err != nil {
		fmt.Fprintf(w, "Error retrieving AES file: %v", err)
		return
	}
	defer aesFile.Close()

	// Ensure filenames match the expected pattern
	if !strings.HasPrefix(aesHandler.Filename, "aes_") || !strings.HasSuffix(handler.Filename, ".png") {
		fmt.Fprintf(w, "Invalid file naming convention. Use aes_<file name>.png and <file name>.png")
		return
	}

	// Extract base filename from the .png file
	baseFileName := strings.TrimSuffix(handler.Filename, ".png")

	// Save the .png file to ./files directory
	filesDir := "./files"
	if _, err := os.Stat(filesDir); os.IsNotExist(err) {
		os.Mkdir(filesDir, os.ModePerm)
	}
	filePath := filepath.Join(filesDir, handler.Filename)
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "Error reading file: %v", err)
		return
	}
	err = ioutil.WriteFile(filePath, fileBytes, 0644)
	if err != nil {
		fmt.Fprintf(w, "Error saving file: %v", err)
		return
	}

	// Save the AES file to ./aes directory
	aesDir := "./aes"
	if _, err := os.Stat(aesDir); os.IsNotExist(err) {
		os.Mkdir(aesDir, os.ModePerm)
	}
	aesFilePath := filepath.Join(aesDir, aesHandler.Filename)
	aesBytes, err := ioutil.ReadAll(aesFile)
	if err != nil {
		fmt.Fprintf(w, "Error reading AES file: %v", err)
		return
	}
	err = ioutil.WriteFile(aesFilePath, aesBytes, 0644)
	if err != nil {
		fmt.Fprintf(w, "Error saving AES file: %v", err)
		return
	}

	// Step 1: Decrypt the AES key
	decryptedAESKey, err := encryptdecrypt.DecryptAES("/mnt/Disk_2/secureTransfer/project/secureTransfer/final", aesBytes)
	if err != nil {
		fmt.Fprintf(w, "Error decrypting AES key: %v", err)
		return
	}

	// Step 2: Decrypt the .png file using the decrypted AES key
	decryptedFileBytes, err := encryptdecrypt.DecodeFile(decryptedAESKey, filePath)
	if err != nil {
		fmt.Fprintf(w, "Error decrypting .png file: %v", err)
		return
	}

	// Step 3: Save the decrypted file to ./info directory
	infoDir := "./info"
	if _, err := os.Stat(infoDir); os.IsNotExist(err) {
		os.Mkdir(infoDir, os.ModePerm)
	}
	decryptedFilePath := filepath.Join(infoDir, "dec_"+baseFileName+".png")
	err = ioutil.WriteFile(decryptedFilePath, decryptedFileBytes, 0644)
	if err != nil {
		fmt.Fprintf(w, "Error saving decrypted file: %v", err)
		return
	}

	fmt.Fprintf(w, "Files uploaded and decrypted successfully. Decrypted file saved as: %s\n", decryptedFilePath)
}

// func main() {
// 	// Simple HTTP server to handle file uploads
// 	http.HandleFunc("/upload", UploadFiles)
// 	fmt.Println("Server started on :8080, use /upload to upload files")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
