package encryptdecrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func encryptFileUsingAES(key []byte, plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

func encodeToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func EncodeFile(aesKey []byte, filePath string) ([]byte, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	encodedFile := encodeToBase64(fileContent)
	encryptedFile, _, err := encryptFileUsingAES(aesKey, []byte(encodedFile))

	if err != nil {
		return nil, err
	}

	return encryptedFile, err

}

// Function to load the RSA public key for encryption
func loadPublicKey(path string) (*rsa.PublicKey, error) {
	// Read the public key file
	pubKeyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	// Parse the public key
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	// Convert to *rsa.PublicKey
	rsaPubKey, ok := pubKey.(ssh.CryptoPublicKey).CryptoPublicKey().(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return rsaPubKey, nil
}

// Function to encrypt the AES key using RSA public key
func encryptWithRSA(publicKey *rsa.PublicKey, aesKey []byte) ([]byte, error) {
	// Encrypt AES key using RSA
	encryptedAESKey, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, aesKey)
	if err != nil {
		return nil, fmt.Errorf("RSA encryption failed: %w", err)
	}
	return encryptedAESKey, nil
}

func EncryptAES(path string, askey []byte) ([]byte, error) {
	rsaPubKey, err := loadPublicKey(path)
	if err != nil {
		return nil, err
	}

	encryptedKey, err := encryptWithRSA(rsaPubKey, askey)
	if err != nil {
		return nil, err
	}

	return encryptedKey, nil
}
