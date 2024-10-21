const chat = document.getElementById('chat');
const messageInput = document.getElementById('message');
const usernameInput = document.getElementById('username');

const token = prompt('Digite seu token JWT');

const socket = new WebSocket(`ws://localhost:7120/ws?token=${token}&room=default`);

socket.onmessage = function(event) {
    const msg = JSON.parse(event.data);
    const messageElement = document.createElement('div');
    messageElement.innerText = `${msg.username}: ${msg.text}`;
    chat.appendChild(messageElement);
};

function sendMessage() {
    const message = {
        username: usernameInput.value,
        text: messageInput.value,
        room: 'default'
    };
    socket.send(JSON.stringify(message));
    messageInput.value = ''; 
}
