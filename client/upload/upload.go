package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"secureTransfer/encryptdecrypt"
)

// Helper function to create a multipart file field
func createFileField(writer *multipart.Writer, fieldname, filename string) error {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create a form file field in the multipart writer
	part, err := writer.CreateFormFile(fieldname, filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("error creating form file field: %v", err)
	}

	// Copy the file content to the form file field
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("error copying file content: %v", err)
	}

	return nil
}

// UploadFiles uploads both the .png file and the aes_ .png file to the server
func UploadFiles(filePath, aesFilePath, url string) error {
	// Prepare a buffer to hold the multipart data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add the main file (<file name>.png) to the request
	if err := createFileField(writer, "file", filePath); err != nil {
		return fmt.Errorf("error adding file field: %v", err)
	}

	// Add the AES-encrypted file (aes_<file name>.png) to the request
	if err := createFileField(writer, "aesFile", aesFilePath); err != nil {
		return fmt.Errorf("error adding AES file field: %v", err)
	}

	// Close the multipart writer to finalize the form data
	err := writer.Close()
	if err != nil {
		return fmt.Errorf("error closing multipart writer: %v", err)
	}

	// Send the POST request to the server
	// Adjust the URL if necessary
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the content type to multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response from the server
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	// Print the server response
	fmt.Println("Server Response:", string(respBody))
	return nil
}

func UploadFilesAutomated(filePath, rsaFilePath string) error {
	// Paths to the files to be uploaded
	aeskey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	url := "http://localhost:8081/upload"
	encryptedFile, err := encryptdecrypt.EncodeFile([]byte(aeskey), filePath)
	if err != nil {
		return err
	}
	name := "encrypted_" + filepath.Base(filePath)
	err = ioutil.WriteFile(name, encryptedFile, 0644)
	if err != nil {
		return err
	}

	nameAES := "aes_" + name

	encryptedKey, err := encryptdecrypt.EncryptAES(rsaFilePath, []byte(aeskey))

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(nameAES, encryptedKey, 0644)

	// Upload the files
	err = UploadFiles(name, nameAES, url)
	if err != nil {
		return err
	} else {
		fmt.Println("Files uploaded successfully")
		return nil
	}
}
