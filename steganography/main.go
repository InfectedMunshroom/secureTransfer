package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"secureTransfer/client"
	"strings"
)

const tempDir = "./tempDir"

var key = []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa") // 32-byte key

func main() {
	http.HandleFunc("/down", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./download.html")
	})
	http.HandleFunc("/decrypt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./deindex.html")
	})
	http.HandleFunc("/encrypt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./main.html")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./index.html") })

	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/upload_encrypt", uploadEncryptHandler)
	http.HandleFunc("/upload_decrypt", uploadDecryptHandler)

	log.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	host := r.FormValue("host")
	file := r.FormValue("file")
	path := r.FormValue("path")

	if host == "" || file == "" || path == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	err := client.DownloadFile(host, file, path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to download file: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File downloaded successfully to %s", path)
}

func uploadEncryptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to read image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	message := r.FormValue("message")
	if message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	filename := r.FormValue("filename")
	if filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".png"
	}

	filename = strings.TrimSuffix(filename, ext) + ext

	imgData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read image", http.StatusInternalServerError)
		return
	}

	url := r.FormValue("serverurl")
	if url == "" {
		http.Error(w, "Server url is required in formt <ip>:<port>", http.StatusBadRequest)
	}

	encryptedMsg, err := encryptAES([]byte(message), key)
	if err != nil {
		http.Error(w, "Unable to encrypt message", http.StatusInternalServerError)
		return
	}

	imgData = append(imgData, []byte("---END---")...)
	imgData = append(imgData, []byte(encryptedMsg)...)

	os.MkdirAll(tempDir, 0755)
	outputPath := filepath.Join(tempDir, filename)
	err = ioutil.WriteFile(outputPath, imgData, 0644)
	if err != nil {
		http.Error(w, "Unable to save image", http.StatusInternalServerError)
		return
	}

	err = client.UploadFilesAutomated(outputPath, "./final.pub", url)
	if err != nil {
		http.Error(w, "Error in uploading image to server", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	fmt.Fprintf(w, "success")
}

func uploadDecryptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to read image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imgData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read image", http.StatusInternalServerError)
		return
	}

	separator := []byte("---END---")
	sepIndex := strings.Index(string(imgData), string(separator))
	if sepIndex == -1 {
		http.Error(w, "No embedded message found in the image", http.StatusBadRequest)
		return
	}

	encryptedMsg := imgData[sepIndex+len(separator):]

	decryptedMsg, err := decryptAES([]byte(encryptedMsg), key)
	if err != nil {
		http.Error(w, "Unable to decrypt message", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Decrypted message: %s", decryptedMsg)
}

func encryptAES(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext), nil
}

func decryptAES(ciphertext []byte, key []byte) (string, error) {
	ciphertextBytes, err := hex.DecodeString(string(ciphertext))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}
