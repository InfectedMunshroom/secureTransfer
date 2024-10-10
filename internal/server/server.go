package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Fprintf(w, "Error parsing form: %v", err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "Error retrieving the file: %v", err)
		return
	}
	defer file.Close()

	destDir := "./file"
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		err = os.Mkdir(destDir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(w, "Error creating directory: %v", err)
			return
		}
	}

	destPath := filepath.Join(destDir, handler.Filename)
	dst, err := os.Create(destPath)
	if err != nil {
		fmt.Fprintf(w, "Error creating the file: %v", err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Fprintf(w, "Error saving the file: %v", err)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %v\n", handler.Filename)
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Filename is missing", http.StatusBadRequest)
		return
	}

	filepath := "./file/" + filename
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	io.Copy(w, file)
}

/*
func main() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)

	fmt.Println("Server started on 10.0.2.3:8080...")
	err := http.ListenAndServe("10.0.2.3:8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
*/
