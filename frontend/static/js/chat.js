function openChat(username) {
    if (!currentUsername) {
        showToast("red", "login to open chat");
        return;
    }
    selectedUser = username;

    // Set User Name
    document.getElementById('chat-username').innerText = username;
    if (typeof setChatStatusByUsername === "function") {
        setChatStatusByUsername(username);
    }

    // Clear previous messages
    const chatMessagesContainer = document.getElementById('chat-messages-container');
    chatMessagesContainer.innerHTML = '';

    initMessagePagination();
    loadMessages(true);

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
    renderMessage(data, { prepend: false, keepScroll: false });
}

function renderMessage(data, { prepend, keepScroll }) {
    const chatMessagesContainer = document.getElementById('chat-messages-container');
    if (!chatMessagesContainer) return;

    const msgDiv = document.createElement('div');

    msgDiv.className = data.from === selectedUser ? 'msg msg-out' : 'msg msg-in';
    msgDiv.innerText = data.msg;

    if (prepend) {
        const prevHeight = chatMessagesContainer.scrollHeight;
        chatMessagesContainer.prepend(msgDiv);
        if (keepScroll) {
            const newHeight = chatMessagesContainer.scrollHeight;
            chatMessagesContainer.scrollTop += newHeight - prevHeight;
        }
    } else {
        chatMessagesContainer.appendChild(msgDiv);
        chatMessagesContainer.scrollTop = chatMessagesContainer.scrollHeight;
    }
}

function sendMessageFromButton() {
    const input = document.getElementById('messageInput');
    if (input && input.value.trim()) {
        sendMessage(input.value);
        input.value = '';
    }
}

const MESSAGES_PAGE_SIZE = 10;
let messagesOffset = 0;
let messagesLoading = false;
let messagesHasMore = true;
let messagesScrollHandler = null;

function initMessagePagination() {
    messagesOffset = 0;
    messagesLoading = false;
    messagesHasMore = true;

    const container = document.getElementById('chat-messages-container');
    if (!container) return;
    if (messagesScrollHandler) {
        container.removeEventListener("scroll", messagesScrollHandler);
    }
    messagesScrollHandler = throttle(() => {
        if (messagesLoading || !messagesHasMore) return;
        if (container.scrollTop <= 30) {
            loadMessages(false);
        }
    }, 300);
    container.addEventListener("scroll", messagesScrollHandler);
}

function loadMessages(reset) {
    if (!selectedUser || messagesLoading) return;
    messagesLoading = true;

    const offset = reset ? 0 : messagesOffset;
    fetch(`/api/conversations/messages?id=${selectedUser}&limit=${MESSAGES_PAGE_SIZE}&offset=${offset}`)
        .then(res => res.json())
        .then(data => {
            if (reset) {
                const container = document.getElementById('chat-messages-container');
                if (container) container.innerHTML = "";
                messagesOffset = 0;
                messagesHasMore = true;
            }

            if (data.messages && data.messages.length > 0) {
                const msgs = data.messages.slice().reverse();
                const container = document.getElementById('chat-messages-container');
                const fragment = document.createDocumentFragment();
                msgs.forEach(msg => {
                    const msgDiv = document.createElement('div');
                    msgDiv.className = msg.from === selectedUser ? 'msg msg-out' : 'msg msg-in';
                    msgDiv.innerText = msg.msg;
                    fragment.appendChild(msgDiv);
                });

                if (container) {
                    if (reset) {
                        container.appendChild(fragment);
                        container.scrollTop = container.scrollHeight;
                    } else {
                        const prevHeight = container.scrollHeight;
                        container.prepend(fragment);
                        const newHeight = container.scrollHeight;
                        container.scrollTop += newHeight - prevHeight;
                    }
                }

                messagesOffset += data.messages.length;
                if (data.messages.length < MESSAGES_PAGE_SIZE) {
                    messagesHasMore = false;
                }
            } else {
                messagesHasMore = false;
            }
        })
        .catch(err => console.error("Failed to load messages:", err))
        .finally(() => {
            messagesLoading = false;
        });
}

function throttle(fn, wait) {
    let last = 0;
    let timeout = null;
    return function (...args) {
        const now = Date.now();
        const remaining = wait - (now - last);
        if (remaining <= 0) {
            if (timeout) {
                clearTimeout(timeout);
                timeout = null;
            }
            last = now;
            fn.apply(this, args);
        } else if (!timeout) {
            timeout = setTimeout(() => {
                last = Date.now();
                timeout = null;
                fn.apply(this, args);
            }, remaining);
        }
    };
}
