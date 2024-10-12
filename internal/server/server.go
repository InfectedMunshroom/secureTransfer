package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"secureTransfer/encryptdecrypt"
	"strings"
)

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

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	// Specify the path of the file in the "info" directory to download.
	fileName := r.URL.Query().Get("file") // E.g., "example.png"
	if fileName == "" {
		http.Error(w, "File not specified.", http.StatusBadRequest)
		return
	}
	infoDir := "./info"
	filePath := filepath.Join(infoDir, fileName)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}

	// Encrypt the file data using AES (Assuming EncryptAES is your encryption function)
	aesKey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" // Replace with your actual AES key generation method
	encryptedData, err := encryptdecrypt.EncodeFile([]byte(aesKey), filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in encrypting file for transit: %v", err), http.StatusInternalServerError)
		return
	}

	// Encrypt the AES key using RSA (Assuming EncryptRSA is your encryption function)
	encryptedAESKey, err := encryptdecrypt.EncryptAES("./final.pub", []byte(aesKey))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in encrypting AES key: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a response object containing the encrypted AES key and the encrypted file data
	response := map[string]interface{}{
		"aes_key":   encryptedAESKey,
		"file_data": encryptedData,
	}

	// Send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("File encrypted and sent successfully:", fileName)
}
