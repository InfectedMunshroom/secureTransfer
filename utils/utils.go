package utils

import (
	"image"
	"image/color"
)

// EmbedMessage embeds a message into the least significant bits of the image.
func EmbedMessage(img image.Image, message []byte) image.Image {
	// Convert the image to RGBA
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImg := image.NewRGBA(bounds)

	// Copy the original image to the new image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}

	// Embed the message into the image
	msgLen := len(message)
	for i := 0; i < msgLen; i++ {
		if i >= width*height {
			break // Prevent writing outside the image bounds
		}
		x := i % width
		y := i / width
		r, g, b, a := newImg.At(x, y).RGBA()

		// Set the least significant bit of the red channel to the message bit
		r = (r & 0xFE) | (uint32(message[i]) & 1) // Modify LSB
		newImg.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
	}

	return newImg
}

// ExtractMessage extracts a message from the least significant bits of the image.
func ExtractMessage(img image.Image) []byte {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	message := make([]byte, width*height)

	for i := 0; i < len(message); i++ {
		x := i % width
		y := i / width
		r, _, _, _ := img.At(x, y).RGBA()

		// Read the LSB of the red channel
		message[i] = byte(r & 1)
	}

	return message
}
