package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"secureTransfer/encryptdecrypt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func UploadFiles(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Fprintf(w, "Error parsing form: %v", err)
		return
	}

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

	if !strings.HasPrefix(aesHandler.Filename, "aes_") || !strings.HasSuffix(handler.Filename, ".png") {
		fmt.Fprintf(w, "Invalid file naming convention. Use aes_<file name>.png and <file name>.png")
		return
	}

	attestation := r.FormValue("attestation")
	if attestation == "" {
		http.Error(w, "Attestation missing", http.StatusBadRequest)
		return
	}

	baseFileName := strings.TrimSuffix(handler.Filename, ".png")

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

	decryptedAESKey, err := encryptdecrypt.DecryptAES("final", aesBytes)
	if err != nil {
		fmt.Fprintf(w, "Error decrypting AES key: %v", err)
		return
	}

	decryptedFileBytes, err := encryptdecrypt.DecodeFile(decryptedAESKey, filePath)
	if err != nil {
		fmt.Fprintf(w, "Error decrypting .png file: %v", err)
		return
	}

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

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		http.Error(w, "Database conection error", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	_, err = db.Exec("INSERT INTO attestations (filename, attestation) VALUES(?,?)", handler.Filename, attestation)
	if err != nil {
		http.Error(w, "Failed to store attestation", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Files uploaded and decrypted successfully. Decrypted file saved as: %s\n", decryptedFilePath)
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File not specified.", http.StatusBadRequest)
		return
	}
	infoDir := "./info"
	filePath := filepath.Join(infoDir, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}

	aesKey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	encryptedData, err := encryptdecrypt.EncodeFile([]byte(aesKey), filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in encrypting file for transit: %v", err), http.StatusInternalServerError)
		return
	}

	encryptedAESKey, err := encryptdecrypt.EncryptAES("final.pub", []byte(aesKey))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in encrypting AES key: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"aes_key":   encryptedAESKey,
		"file_data": encryptedData,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("File encrypted and sent successfully:", fileName)
}
