document.getElementById('uploadForm').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const formData = new FormData(this);
    const loading = document.getElementById('loading');
    
    loading.style.display = 'block';

    fetch('/upload_encrypt', {
        method: 'POST',
        body: formData
    })
    .then(response => response.text())
    .then(data => {
        showNotification('success', 'File uploaded and encrypted successfully!');
        setTimeout(() => {
            window.location.href = '/';
        }, 2000);
    })
    .catch(error => {
        showNotification('error', 'An error occurred: ' + error.message);
    })
    .finally(() => {
        loading.style.display = 'none';
    });
});

function showNotification(type, message) {
    const notification = document.getElementById('notification');
    notification.className = 'notification ' + (type === 'success' ? 'success' : 'error');
    notification.textContent = message;
    notification.style.display = 'block';
    setTimeout(() => {
        notification.style.display = 'none';
    }, 3000);
}
