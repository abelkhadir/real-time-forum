// WebSocket and Chat Functionality
let selectedUser;
let ws;

const counter = document.querySelector(".notifications-counter");

// Initialize WebSocket connection
function initWebSocket() {
    ws = new WebSocket(`ws://${window.location.host}/ws`);

    ws.onopen = () => {
        console.log("WebSocket connected");
    };

    ws.onmessage = (e) => {
        const data = JSON.parse(e.data);

        if (data.type === "UpdateNotifs") {
            // add backend
            updateNotifications(data);
            return;
        }

        if (data.type === "UpdatePosts") {
            // todo: load posts using ws (ma3rftch wach required, khli tanchofo)
            return;
        }

        if (data.type === "UpdateContacts") {
            loadContacts(data.contacts, data.username);
            return;
        }

        if (data.type === "UpdateMessages") {
            console.log("got", data.type);
            displayMessage(data);
            return
        }
        return
    }

    ws.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    ws.onclose = () => {
        console.log("WebSocket disconnected");
    };
};



function updateNotifications() {
    const container = document.getElementById("notifications-container");

    // example payload
    const data = {
        count: 3,
        notifications: [
            "New message from abde",
            "New message from abde2",
            "New message from abde3"
        ]
    };
    
    // counter
    if (data.count > 0) {
        counter.textContent = data.count;
        counter.classList.remove("hidden");
    } else {
        counter.classList.add("hidden");
    }

    // DOM list
    container.innerHTML = ""; // reset
    data.notifications.forEach(msg => {
        const div = document.createElement("div");
        div.className = "notification-item";
        div.textContent = msg;
        container.appendChild(div);
    });
}