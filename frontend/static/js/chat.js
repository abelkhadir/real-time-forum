let selectedUser;
function openChat(username) {
    selectedUser = username

    // Set User Name
    document.getElementById('chat-username').innerText = username;

    // Switch Sidebar Views
    document.getElementById('friends-list').classList.add('hidden');
    document.getElementById('chat-conversation').classList.remove('hidden');
}

function closeChat() {
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
        e.preventDefault()
        const msg = e.target.value

        if (msg != "") {
            sendMessage(msg)
        }

        e.target.value = ""
    }
}

function sendMessage(msg) {
    if (ws.readyState === WebSocket.OPEN) {
        console.log("sending", msg);
        ws.send(msg);
    } else {
        console.log("WebSocket not open, cannot send yet");
    }
}

const ws = new WebSocket(`ws://${window.location.host}/ws`);

ws.onopen = () => {
    console.log("connected");
};

ws.onmessage = (e) => {
    const data = JSON.parse(e.data);
    console.log("Message from", data.from, ":", data.msg);
};

// send a message to a specific client
function sendMessage(msg) {
    const payload = {
        to: selectedUser,
        msg: msg,
    };
    ws.send(JSON.stringify(payload));
}
