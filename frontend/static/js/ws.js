// WebSocket and Chat Functionality
let selectedUser;
let ws;

const counter = document.querySelector(".notifications-counter");
let notifications = [];
let unreadCount = 0;

function renderNotifications() {
    const container = document.getElementById("notifications-container");
    if (!container) return;

    container.innerHTML = "";
    if (notifications.length === 0) {
        const empty = document.createElement("div");
        empty.className = "notification-item";
        empty.textContent = "No new notifications";
        container.appendChild(empty);
        return;
    }

    notifications.forEach((n) => {
        const div = document.createElement("div");
        div.className = "notification-item";
        div.textContent = n.text;
        div.addEventListener("click", () => {
            openChat(n.from);
            const notifMenu = document.getElementById("notif-menu");
            if (notifMenu) notifMenu.classList.add("hidden");
        });
        container.appendChild(div);
    });
}

function updateCounter() {
    if (!counter) return;
    if (unreadCount > 0) {
        counter.textContent = String(unreadCount);
        counter.classList.remove("hidden");
    } else {
        counter.textContent = "0";
        counter.classList.add("hidden");
    }
}

function addNotification(from, msg) {
    const preview = msg && msg.length > 80 ? `${msg.slice(0, 77)}...` : (msg || "");
    const text = preview ? `Message from ${from}: ${preview}` : `Message from ${from}`;
    notifications.unshift({ from, msg, text });
    if (notifications.length > 50) notifications.length = 50;
    unreadCount += 1;
    updateCounter();
    renderNotifications();
}

function markNotificationsRead() {
    fetch("/api/notifications/read", { method: "POST" }).catch(() => {});
    unreadCount = 0;
    updateCounter();
}

function fetchNotifications() {
    return fetch("/api/notifications")
        .then((res) => {
            if (!res.ok) return null;
            return res.json();
        })
        .then((data) => {
            if (!data || !Array.isArray(data.notifications)) return;
            notifications.length = 0;
            data.notifications.forEach((n) => {
                const preview = n.msg && n.msg.length > 80 ? `${n.msg.slice(0, 77)}...` : (n.msg || "");
                const text = preview ? `New message from ${n.from}: ${preview}` : `New message from ${n.from}`;
                notifications.push({ from: n.from, msg: n.msg, text });
            });
            unreadCount = typeof data.count === "number" ? data.count : notifications.length;
            updateCounter();
            renderNotifications();
        })
        .catch(() => {});
}

// Initialize WebSocket connection
function initWebSocket() {
    ws = new WebSocket(`ws://${window.location.host}/ws`);

    ws.onopen = () => {};

    ws.onmessage = (e) => {
        const data = JSON.parse(e.data);

        if (data.type === "UpdateNotifs") {
            // Reserved for future backend-driven notifications
            return;
        }

        if (data.type === "UpdateMessages") {
            if (selectedUser && data.from === selectedUser) {
                displayMessage(data);
            } else {
                addNotification(data.from, data.msg);
            }
            return;
        }

        if (data.type === "UpdatePosts") {
            if (data.post) {
                addPostToFeed(data.post);
            }
            return;
        }

        if (data.type === "UpdateContacts") {
            loadContacts(data.contacts, data.username);
            return;
        }

        return
    }

    ws.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    ws.onclose = () => {};
};

updateCounter();
