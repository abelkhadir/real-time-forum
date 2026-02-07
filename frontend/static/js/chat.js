
function openChat(username) {
    selectedUser = username;

    // Set User Name
    document.getElementById('chat-username').innerText = username;
    if (typeof setChatStatusByUsername === "function") {
        setChatStatusByUsername(username);
    }

    // Clear previous messages
    const chatMessagesContainer = document.getElementById('chat-messages-container');
    chatMessagesContainer.innerHTML = '';

    prevMessages(username);

    // Switch Sidebar Views
    document.getElementById('friends-list').classList.add('hidden');
    document.getElementById('chat-conversation').classList.remove('hidden');
}

function closeChat() {
    selectedUser = null;
    document.getElementById('chat-conversation').classList.add('hidden');
    document.getElementById('friends-list').classList.remove('hidden');
}

function toggleMobileSidebar() {
    const sidebar = document.getElementById('right-sidebar');
    if (sidebar.style.display === 'flex') {
        sidebar.style.display = 'none';
    } else {
        sidebar.style.display = 'flex';
        sidebar.classList.add('active');
    }
}

function readMessage(e) {
    if (e.key === "Enter") {
        e.preventDefault();
        const msg = e.target.value;

        if (msg != "") {
            sendMessage(msg);
        }

        e.target.value = "";
    }
}

function sendMessage(msg) {
    if (!selectedUser) {
        showToast("red", "No user selected");
        return;
    }

    if (ws && ws.readyState === WebSocket.OPEN) {
        const payload = {
            to: selectedUser,
            msg: msg,
        };
        ws.send(JSON.stringify(payload));
        displayMessage(payload);
    } else {
        showToast("red", "WebSocket not connected");
    }
}

function displayMessage(data) {
    const chatMessagesContainer = document.getElementById('chat-messages-container');
    if (!chatMessagesContainer) return;

    const msgDiv = document.createElement('div');
    
    msgDiv.className = data.from === selectedUser ? 'msg msg-out' : 'msg msg-in';
    msgDiv.innerText = data.msg;

    chatMessagesContainer.appendChild(msgDiv);
    chatMessagesContainer.scrollTop = chatMessagesContainer.scrollHeight;
}

function sendMessageFromButton() {
    const input = document.getElementById('messageInput');
    if (input && input.value.trim()) {
        sendMessage(input.value);
        input.value = '';
    }
}

function prevMessages(id) {
    fetch(`/api/conversations/messages?id=${id}&limit=50`)
        .then(res => res.json())
        .then(data => {
            if (data.messages && data.messages.length > 0) {
                data.messages.forEach(msg => displayMessage(msg));
            }
        })
        .catch(err => console.error("Failed to load messages:", err));
}
