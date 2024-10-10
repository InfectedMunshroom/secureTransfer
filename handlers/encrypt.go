package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hello/utils"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
)

// EncryptResponse struct for JSON response
type EncryptResponse struct {
	Key      string `json:"key"`
	ImageURL string `json:"imageURL"`
}

// EncryptHandler handles the encryption of a message embedded in an image.
func EncryptHandler(w http.ResponseWriter, r *http.Request) {
	// Limit the form size to 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	message := r.FormValue("message")
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to read image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Decode the uploaded image
	img, err := png.Decode(file)
	if err != nil {
		http.Error(w, "Unable to decode image", http.StatusInternalServerError)
		return
	}

	// Generate a random key for encryption
	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		http.Error(w, "Error generating key", http.StatusInternalServerError)
		return
	}

	// Encrypt the message
	encryptedMessage, err := utils.EncryptAES(key, []byte(message))
	if err != nil {
		http.Error(w, "Error encrypting message", http.StatusInternalServerError)
		return
	}

	// Ensure EmbedMessage function is implemented in utils
	stegoImage := utils.EmbedMessage(img, encryptedMessage)

	// Save the stego image to a file
	stegoImagePath := filepath.Join("stego_images", "stego_image.png")
	stegoFile, err := os.Create(stegoImagePath)
	if err != nil {
		http.Error(w, "Unable to create stego image file", http.StatusInternalServerError)
		return
	}
	defer stegoFile.Close()

	// Encode and save the stego image
	if err := png.Encode(stegoFile, stegoImage); err != nil {
		http.Error(w, "Unable to encode image", http.StatusInternalServerError)
		return
	}

	// Respond with JSON containing the encryption key and image URL
	response := EncryptResponse{
		Key:      hex.EncodeToString(key),
		ImageURL: fmt.Sprintf("http://localhost:8080/%s", stegoImagePath), // Change to your server's address
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
