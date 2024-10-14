package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"secureTransfer/encryptdecrypt"
)

type DownloadResponse struct {
	AESKey   []byte `json:"aes_key"`
	FileData []byte `json:"file_data"`
}

func DownloadFile(host, file, downloadPath string) error {
	// url := "http://localhost:8081/download?file=dec_encrypted_output.png"
	url := "http://" + host + "/download?file=" + file

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error during HTTP request:%v", err)

	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("Server returned status %d: %s\n", response.StatusCode, string(body))
	}

	var downloadResponse DownloadResponse
	err = json.NewDecoder(response.Body).Decode(&downloadResponse)
	if err != nil {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Error decoding JSON response:", err)
		fmt.Println("Raw response body:", string(body))
		return fmt.Errorf("Error decoding JSON Response: %v", err)
	}

	decryptedAESKey, err := encryptdecrypt.DecryptAES("./../final", downloadResponse.AESKey)
	if err != nil {
		return fmt.Errorf("Error in decoding the AES key: %v", err)

	}
	fmt.Println("Decrypted AES key:", string(decryptedAESKey))

	tempFilePath := "temp"
	err = ioutil.WriteFile(tempFilePath, downloadResponse.FileData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing encrypted data to 'temp':%v", err)

	}
	fmt.Println("Encrypted file data written to 'temp'.")

	decryptedFileData, err := encryptdecrypt.DecodeFile(decryptedAESKey, tempFilePath)
	if err != nil {
		return fmt.Errorf("Error decrypting file: %v", err)

	}

	cleanPath := filepath.Clean(downloadPath)
	err = ioutil.WriteFile(cleanPath, decryptedFileData, 0644)
	if err != nil {
		return fmt.Errorf("Error saving the decrypted file: %v", err)

	}

	fmt.Println("File downloaded and decrypted successfully at:", cleanPath)
	return nil
}
