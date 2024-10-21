let websocket;
let username = '';
const userColors = {}; 

function getColorForUser(username) {
    if (userColors[username]) {
        return userColors[username];
    }

    const color = `#${Math.floor(Math.random()*16777215).toString(16)}`;
    userColors[username] = color;  
    return color;
}

function createRoom() {
    const roomName = document.getElementById('room-name').value;
    const password = document.getElementById('room-password').value;
    username = document.getElementById('user-name').value || 'Usuário';

    fetch('http://localhost:7120/create-room', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name: roomName, password: password })
    })
    .then(response => response.json())
    .then(data => {
        alert('Sala criada com sucesso: ' + data.id);
    })
    .catch(error => {
        console.error('Erro:', error);
    });
}

function joinRoom() {
    const roomId = document.getElementById('join-room-id').value;
    const password = document.getElementById('join-room-password').value;
    username = document.getElementById('join-user-name').value || 'Usuário';
         fetch('http://localhost:7120/join-room', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ room_id: roomId, password: password })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Erro ao entrar na sala: ' + response.statusText);
        }
        return response.json();
    })
    .then(data => {
        alert(data.message);
        document.getElementById('room-creation').style.display = 'none';
        document.getElementById('room-join').style.display = 'none';
        document.getElementById('current-room-name').innerText = roomId;
        document.getElementById('chat').style.display = 'block';
        connectWebSocket(roomId);
    })
    .catch(error => {
        console.error('Erro:', error);
    });
}

function connectWebSocket(roomId) {
    websocket = new WebSocket(`ws://localhost:7120/ws?room=${roomId}`);

    websocket.onmessage = function(event) {
        const msg = JSON.parse(event.data);
        displayMessage(msg);
    };

    websocket.onclose = function() {
        console.log('WebSocket closed');
    };
}

// function sendMessage() {
//     const messageInput = document.getElementById('message-input');
//     const messageText = messageInput.value;

//     const message = {
//         username: 'User',  
//         text: messageText,
//         room: document.getElementById('current-room-name').innerText
//     };

//     websocket.send(JSON.stringify(message));
//     messageInput.value = '';
// }
function sendMessage() {
    const messageInput = document.getElementById('message-input');
    const messageText = messageInput.value;

    if (messageText.trim() === '') return;  

    const message = {
        username: username,
        text: messageText,
        room: document.getElementById('current-room-name').innerText
    };

    websocket.send(JSON.stringify(message));
    messageInput.value = '';   
}
function sendFile() {
    const fileInput = document.getElementById('file-upload');
    const file = fileInput.files[0];
    
    if (file) {
        const formData = new FormData();
        formData.append('file', file);

        fetch('http://localhost:7120/upload', {
            method: 'POST',
            body: formData
        })
        .then(response => response.text())  
        .then(fileURL => {
            const message = {
                username: username,
                file_url: fileURL,  
                room: document.getElementById('current-room-name').innerText
            };
        
            websocket.send(JSON.stringify(message));
        })
        .catch(error => console.error('Erro ao enviar arquivo:', error));
    }
}

websocket.onmessage = function(event) {
    const msg = JSON.parse(event.data);
    displayMessage(msg);
};

function displayMessage(msg) {
    const messagesDiv = document.getElementById('messages');
    const messageElement = document.createElement('div');

    messageElement.classList.add('mb-2', 'p-2', 'rounded-lg', 'shadow', 'bg-gray-800', 'text-gray-200'); 

    if (msg.text) {
        messageElement.innerText = `${msg.username}: ${msg.text}`;
    } else if (msg.file_url) {
        const imageExtensions = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp'];
        const audioExtensions = ['mp3', 'wav', 'ogg', 'm4a'];
        const videoExtensions = ['mp4', 'webm', 'ogg'];
        const fileExtension = msg.file_url.split('.').pop().toLowerCase();

        if (imageExtensions.includes(fileExtension)) {
            const img = document.createElement('img');
            img.src = msg.file_url;
            img.alt = 'Imagem enviada';
            img.classList.add('max-w-xs', 'rounded', 'shadow', 'border', 'border-gray-600');
            messageElement.appendChild(img);
        } else if (audioExtensions.includes(fileExtension)) {
            const audio = document.createElement('audio');
            audio.controls = true;  
            audio.src = msg.file_url;
            messageElement.appendChild(audio);
        } else if (videoExtensions.includes(fileExtension)) {
            const video = document.createElement('video');
            video.controls = true; 
            video.src = msg.file_url;
            video.classList.add('max-w-xs', 'rounded', 'shadow', 'border', 'border-gray-600');
            messageElement.appendChild(video);
        } else {
            const fileLink = document.createElement('a');
            fileLink.href = msg.file_url;
            fileLink.innerText = 'Clique para ver o arquivo';
            fileLink.target = '_blank';  
            fileLink.classList.add('text-blue-400', 'hover:underline'); 
            messageElement.appendChild(fileLink);
        }
    }

    const usernameElement = document.createElement('span');
    usernameElement.innerText = `${msg.username}: `;
    usernameElement.classList.add('font-semibold', 'text-blue-300');  // Texto do nome de usuário em azul claro
    messageElement.prepend(usernameElement);  

    messagesDiv.appendChild(messageElement);
    messagesDiv.scrollTop = messagesDiv.scrollHeight; 
}




// function displayMessage(msg) {
//     const messagesDiv = document.getElementById('messages');
//     const messageElement = document.createElement('div');
//     messageElement.innerText = `${msg.username}: ${msg.text}`;
//     messagesDiv.appendChild(messageElement);
//     messagesDiv.scrollTop = messagesDiv.scrollHeight; 
// }
