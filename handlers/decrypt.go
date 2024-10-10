package handlers

import (
	"fmt"
	"hello/utils"
	"image/png"
	"net/http"
)

// DecryptHandler handles the extraction of a message from an image.
func DecryptHandler(w http.ResponseWriter, r *http.Request) {
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

	// Extract the message from the image
	encryptedMessage := utils.ExtractMessage(img)

	// Here, you might want to decrypt the message
	// For demonstration purposes, we'll just respond with the extracted message
	fmt.Fprintf(w, "Extracted Message: %s", encryptedMessage)
}
