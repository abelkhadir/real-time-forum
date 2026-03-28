// WebSocket and Chat Functionality
let selectedUser;
let ws;

const counter = document.querySelector(".notifications-counter");
let notifications = [];
let unreadCount = 0;

function formatNotificationText(n) {
    const preview = n.msg && n.msg.length > 80 ? `${n.msg.slice(0, 77)}...` : (n.msg || "");
    if (preview) return `${n.from} (${n.count}) • ${preview}`;
    return `${n.from} (${n.count})`;
}

// renderNotifications draws the current notification list.
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
        div.textContent = formatNotificationText(n);
        div.addEventListener("click", () => {
            openChat(n.from);
            const notifMenu = document.getElementById("notif-menu");
            if (notifMenu) notifMenu.classList.add("hidden");
        });
        container.appendChild(div);
    });
}

// updateCounter refreshes the unread notification badge.
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

// addNotification adds a new unread notification to the list.
function addNotification(from, msg) {
    const notifMenu = document.getElementById("notif-menu");
    const isNotifMenuOpen = notifMenu && !notifMenu.classList.contains("hidden");
    const existingIdx = notifications.findIndex((n) => n.from === from);
    if (existingIdx >= 0) {
        const existing = notifications[existingIdx];
        existing.count += 1;
        existing.msg = msg || existing.msg;
        notifications.splice(existingIdx, 1);
        notifications.unshift(existing);
    } else {
        notifications.unshift({ from, msg: msg || "", count: 1 });
    }
    if (notifications.length > 50) notifications.length = 50;
    if (!isNotifMenuOpen) unreadCount += 1;
    updateCounter();
    renderNotifications();
}

// markNotificationsRead clears unread notifications on the server and UI.
function markNotificationsRead(from = "") {
    const url = from
        ? `/api/notifications/read?from=${encodeURIComponent(from)}`
        : "/api/notifications/read";
    fetch(url, { method: "POST" }).catch(() => {});
    if (from) {
        notifications = notifications.filter((n) => n.from !== from);
        unreadCount = notifications.reduce((sum, n) => sum + n.count, 0);
    } else {
        unreadCount = 0;
    }
    updateCounter();
    renderNotifications();
}

// fetchNotifications loads unread notifications from the backend.
function fetchNotifications() {
    return fetch("/api/notifications")
        .then((res) => {
            if (!res.ok) return null;
            return res.json();
        })
        .then((data) => {
            if (!data || !Array.isArray(data.notifications)) return;
            notifications.length = 0;
            const grouped = new Map();
            data.notifications.forEach((n) => {
                const existing = grouped.get(n.from);
                if (existing) {
                    existing.count += 1;
                } else {
                    grouped.set(n.from, { from: n.from, msg: n.msg || "", count: 1 });
                }
            });
            notifications.push(...grouped.values());
            unreadCount = typeof data.count === "number"
                ? data.count
                : notifications.reduce((sum, n) => sum + n.count, 0);
            updateCounter();
            renderNotifications();
        })
        .catch(() => {});
}

// initWebSocket connects the realtime event stream and dispatches updates.
function initWebSocket() {
    ws = new WebSocket(`ws://${window.location.host}/ws`);

    ws.onopen = () => {};

    ws.onmessage = (e) => {
        const data = JSON.parse(e.data);

        if (data.type === "UpdateMessages") {
            const isActiveConversation = selectedUser && (data.from === selectedUser || data.to === selectedUser);

            if (isActiveConversation) {
                displayMessage(data);
                if (data.from === selectedUser) {
                    markNotificationsRead(selectedUser);
                }
            } else {
                const selfUsername = typeof currentUsername === "string" ? currentUsername : "";
                const isOwnMessage = selfUsername !== "" && data.from === selfUsername;
                if (!isOwnMessage) {
                    addNotification(data.from, data.msg);
                }
            }
            return;
        }

        if (data.type === "UpdatePosts") {
            if (typeof getPosts === "function") {
                getPosts();
            }
            return;
        }

        if (data.type === "UpdateComments") {
            if (typeof getPosts === "function") {
                getPosts();
            }
            const postView = document.getElementById("post-view");
            if (postView && !postView.classList.contains("hidden") && typeof refreshCurrentPost === "function") {
                refreshCurrentPost();
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
