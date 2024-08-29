document.addEventListener('DOMContentLoaded', () => {
    fetch('/api/user/1')
        .then(response => response.json())
        .then(data => {
            const userDiv = document.getElementById('user');
            userDiv.innerHTML = `ID: ${data.id}, Name: ${data.name}, Age: ${data.age}`;
        })
        .catch(error => console.error('Error:', error));
});
