package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"secureTransfer/internal/server"
)

// Struct to handle login credentials
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Struct for the login response
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Hardcoded credentials
var validUsername = "admin"
var validPassword = "password"

// Handler for serving static files
func staticFileHandler() http.Handler {
	fs := http.FileServer(http.Dir("./static"))
	return http.StripPrefix("/static/", fs)
}

// Handler for login functionality
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials

	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate the credentials
	if creds.Username == validUsername && creds.Password == validPassword {
		response := Response{
			Success: true,
			Message: "Login successful",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		response := Response{
			Success: false,
			Message: "Invalid username or password",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// Redirect handler for root URL ("/")
func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/login.html", http.StatusSeeOther) // Redirect to login page
}

func main() {
	// Serve static files (HTML, CSS, JS)
	http.Handle("/static/", staticFileHandler())

	// Redirect root URL to the login page
	http.HandleFunc("/", redirectToLogin)
	http.HandleFunc("/upload", server.UploadFile)
	http.HandleFunc("/download", server.DownloadFile)

	// Handle login API
	http.HandleFunc("/api/login", loginHandler)

	// Start the server on 10.0.2.3:8080
	fmt.Println("Server is running on http://10.0.2.15:8080")
	err := http.ListenAndServe("10.0.2.15:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
