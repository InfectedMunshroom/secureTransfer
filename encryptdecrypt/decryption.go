package encryptdecrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// Here lies the code to decrypt the file using a known AES Key
func decryptAES(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {

		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println("Error here 3")

		return nil, err
	}

	return plaintext, nil
}

// Base64 decode the content
func decodeFromBase64(encodedData string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedData)
}

func DecodeFile(aesKey []byte, filePath string) ([]byte, error) {
	encryptedFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}
	plaintext, err := decryptAES(aesKey, encryptedFile)
	if err != nil {
		return nil, err
	}

	file, err := decodeFromBase64(string(plaintext))

	return file, err

}

// Here lies the code to decrypt the AES key after transfer
func loadOpenSSHPrivateKey(path string) (*rsa.PrivateKey, error) {
	// Read the OpenSSH private key file
	privateKeyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	// Parse the OpenSSH private key, remember must be in PEM format
	privateKey, err := ssh.ParseRawPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenSSH private key: %w", err)
	}

	// Convert to *rsa.PrivateKey
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}

	return rsaPrivateKey, nil
}

// Function to decrypt the AES key using RSA private key
func decryptWithRSA(privateKey *rsa.PrivateKey, encryptedAESKey []byte) ([]byte, error) {
	// Decrypt AES key using RSA
	decryptedAESKey, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedAESKey)
	if err != nil {
		return nil, fmt.Errorf("RSA decryption failed: %w", err)
	}
	return decryptedAESKey, nil
}

func DecryptAES(path string, encryptedAESKey []byte) ([]byte, error) {
	rsaPrivateKey, err := loadOpenSSHPrivateKey(path)
	if err != nil {
		return nil, fmt.Errorf("Error: Could not parse ssh key.\n", err)
	}

	// encryptedAESKey, err := ioutil.ReadFile(filepath)
	// if err != nil {
	// 	return nil, err
	// }

	aeskey, err := decryptWithRSA(rsaPrivateKey, encryptedAESKey)
	if err != nil {
		return nil, fmt.Errorf("Error: Could not decrypt.\n", err)
	}
	return aeskey, nil
}
