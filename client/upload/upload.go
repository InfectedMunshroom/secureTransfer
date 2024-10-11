package upload

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"secureTransfer/encryptdecrypt"
)

func UploadFile(url, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("could not create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("could not close writer: %v", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response: %v", err)
	}

	fmt.Println("Server response:", string(respBody))
	return nil
}

func main() {
	// Check if enough arguments are passed
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run client.go <https://server_url/upload> <file_path>")
		return
	}

	// Get URL and filename from command-line arguments
	url := os.Args[1]
	filename := os.Args[2]
	aeskey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	encryptedFile, err := encryptdecrypt.EncodeFile([]byte(aeskey), filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	var name string = "encrypted_" + filename
	err = ioutil.WriteFile(name, encryptedFile, 0644)

	// Call the upload function
	err = UploadFile(url, "temp")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("File uploaded successfully!")
	}
}
