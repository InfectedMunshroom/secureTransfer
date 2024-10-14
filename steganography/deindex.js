document.getElementById("decryptForm").addEventListener("submit", function(event) {
    event.preventDefault();
    const formData = new FormData(this);
    fetch("/upload_decrypt", {
        method: "POST",
        body: formData,
    })
    .then(response => {
        if (!response.ok) {
            throw new Error("Network response was not ok " + response.statusText);
        }
        return response.text();
    })
    .then(data => {
        showNotification(data, true);
    })
    .catch(error => {
        showNotification("Error: " + error.message, false);
    });
});

function showNotification(message, isSuccess) {
    const notification = document.getElementById("notification");
    notification.style.display = "block";
    notification.className = "notification " + (isSuccess ? "success" : "error");
    notification.innerText = message;
}
