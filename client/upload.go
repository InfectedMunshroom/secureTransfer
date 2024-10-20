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
	"time"
)

func createFileField(writer *multipart.Writer, fieldname, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile(fieldname, filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("error creating form file field: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("error copying file content: %v", err)
	}

	return nil
}

func createField(writer *multipart.Writer, fieldname, value string) error {
	part, err := writer.CreateFormField(fieldname)
	if err != nil {
		return fmt.Errorf("error creating form field: %v", err)
	}
	_, err = part.Write([]byte(value))
	if err != nil {
		return fmt.Errorf("error writing form field: %v", err)
	}
	return nil
}

func UploadFiles(filePath, aesFilePath, attestation, url string) error {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	if err := createFileField(writer, "file", filePath); err != nil {
		return fmt.Errorf("error adding file field: %v", err)
	}

	if err := createFileField(writer, "aesFile", aesFilePath); err != nil {
		return fmt.Errorf("error adding AES file field: %v", err)
	}

	if err := createField(writer, "attestation", attestation); err != nil {
		return fmt.Errorf("error adding attestation field: %v", err)
	}

	err := writer.Close()
	if err != nil {
		return fmt.Errorf("error closing multipart writer: %v", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	fmt.Println("Server Response:", string(respBody))
	return nil
}

func UploadFilesAutomated(filePath, rsaFilePath, url string) error {
	aeskey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	url = "http://" + url + "/upload"
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
	attestation := fmt.Sprintf("Attestation-%d", time.Now().Unix())
	err = UploadFiles(name, nameAES, attestation, url)
	if err != nil {
		return err
	} else {
		fmt.Println("Files uploaded successfully")
		return nil
	}
}
