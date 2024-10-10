// Encrypt Function
async function encrypt() {
    const message = document.getElementById("encryption-text").value;
    const imageFile = document.getElementById("encryption-file").files[0];
    const resultDiv = document.getElementById("encryption-result");

    if (!message || !imageFile) {
        resultDiv.innerHTML = "Please provide both a message and an image file.";
        return;
    }

    // Create a FormData object to send to the server
    const formData = new FormData();
    formData.append("message", message);
    formData.append("image", imageFile);

    try {
        // Send POST request to the encryption endpoint
        const response = await fetch('http://localhost:8080/encrypt', {
            method: 'POST',
            body: formData
        });

        // Check if the request was successful
        if (!response.ok) {
            throw new Error('Encryption failed');
        }

        // Read the response as JSON
        const result = await response.json();

        // Create a link to download the stego image
        const link = document.createElement('a');
        link.href = result.imageURL; // Assuming you serve the image correctly
        link.download = 'stego_image.png';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link); // Clean up the link

        // Display the encryption key in the result div
        resultDiv.innerHTML = `Encryption Key: ${result.key}`;
    } catch (error) {
        resultDiv.innerHTML = `Error: ${error.message}`;
    }
}
