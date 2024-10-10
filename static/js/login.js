function login() {
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username: username, password: password })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            window.location.href = "/static/success.html"; // Redirect to success page
        } else {
            document.getElementById("errorMessage").innerText = data.message; // Show error message
        }
    })
    .catch(error => console.error('Error:', error));
}

