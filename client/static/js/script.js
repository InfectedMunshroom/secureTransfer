document.getElementById("uploadForm").addEventListener("submit", function(event) {
    event.preventDefault(); // Prevent the form from submitting normally

    const filePath = document.getElementById("filePath").value; // Get the file path from input

    // Send the file path to the server using Fetch API
    fetch("/upload", {
        method: "POST",
        headers: {
            "Content-Type": "application/json" // Send as JSON
        },
        body: JSON.stringify({ filename: filePath }) // Convert to JSON
    })
    .then(response => {
        if (!response.ok) {
            throw new Error("File upload failed");
        }
        return response.text();
    })
    .then(data => {
        document.getElementById("response").innerText = data; // Show the response
    })
    .catch(error => {
        document.getElementById("response").innerText = error.message; // Handle errors
    });
});
