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
	"strings"
)

const tempDir = "./temp"

var key = []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa") // 32-byte key

func main() {
	// Serve the index.html file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html") // Serve index.html from the root directory
	})

	http.HandleFunc("/upload_encrypt", uploadEncryptHandler)
	http.HandleFunc("/upload_decrypt", uploadDecryptHandler)

	log.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
func uploadEncryptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read uploaded image
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

	imgData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read image", http.StatusInternalServerError)
		return
	}

	// Encrypt the message
	encryptedMsg, err := encryptAES([]byte(message), key)
	if err != nil {
		http.Error(w, "Unable to encrypt message", http.StatusInternalServerError)
		return
	}

	// Append the encrypted message with a separator to the image data
	imgData = append(imgData, []byte("---END---")...)
	imgData = append(imgData, []byte(encryptedMsg)...)

	// Save the modified image to the /temp directory
	os.Mkdir(tempDir, 0755)
	outputPath := filepath.Join(tempDir, "output.png")
	err = ioutil.WriteFile(outputPath, imgData, 0644)
	if err != nil {
		http.Error(w, "Unable to save image", http.StatusInternalServerError)
		return
	}

	// Return a simple success message for the frontend to catch
	fmt.Fprintf(w, "success")
}

func uploadDecryptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read uploaded image
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

	// Look for the "---END---" separator in binary data
	separator := []byte("---END---")
	sepIndex := strings.Index(string(imgData), string(separator))
	if sepIndex == -1 {
		http.Error(w, "No embedded message found in the image", http.StatusBadRequest)
		return
	}

	// Extract the encrypted message after the separator
	encryptedMsg := imgData[sepIndex+len(separator):]

	// Decrypt the message
	decryptedMsg, err := decryptAES([]byte(encryptedMsg), key)
	if err != nil {
		http.Error(w, "Unable to decrypt message", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Decrypted message: %s", decryptedMsg)
}

// encryptAES encrypts plaintext using AES and returns the encrypted message in hex format
func encryptAES(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Generate a random IV
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext), nil
}

// decryptAES decrypts a hex-encoded AES ciphertext and returns the plaintext message
func decryptAES(ciphertext []byte, key []byte) (string, error) {
	// Convert from hex encoding
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
