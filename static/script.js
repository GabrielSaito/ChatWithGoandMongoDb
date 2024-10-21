let websocket;

function createRoom() {
    const roomName = document.getElementById('room-name').value;
    const password = document.getElementById('room-password').value;

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

function sendMessage() {
    const messageInput = document.getElementById('message-input');
    const messageText = messageInput.value;

    const message = {
        username: 'User',  
        text: messageText,
        room: document.getElementById('current-room-name').innerText
    };

    websocket.send(JSON.stringify(message));
    messageInput.value = '';
}

function displayMessage(msg) {
    const messagesDiv = document.getElementById('messages');
    const messageElement = document.createElement('div');
    messageElement.innerText = `${msg.username}: ${msg.text}`;
    messagesDiv.appendChild(messageElement);
    messagesDiv.scrollTop = messagesDiv.scrollHeight; 
}
