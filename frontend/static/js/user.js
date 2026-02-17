let contactsByName = {};
let currentUsername = "";
let lastContacts = null;

function setChatStatusByUsername(username) {
    const statusEl = document.getElementById("chat-status");
    if (!statusEl) return;
    const contact = contactsByName[username];
    const online = contact ? contact.Online : false;
    statusEl.textContent = online ? "Online" : "Offline";
    statusEl.classList.toggle("online", online);
    statusEl.classList.toggle("offline", !online);
}

function loadContacts(contacts) {
    if (!contacts || contacts.length === 0) {
        return;
    }

    lastContacts = contacts;
    const selfUsername = currentUsername || username || "";

    let posts = 0

    const div = document.getElementById("friends-list");
    div.innerHTML = ""; //  clear old list
    contactsByName = {};

    contacts.forEach(contact => {
        contactsByName[contact.Username] = contact;
        if (selfUsername && contact.Username === selfUsername) return;
        posts++;
        const friend = document.createElement("div");
        friend.className = "friends-item";

        const item = document.createElement("div");
        item.className = "friend-item";
        item.addEventListener("click", () => openChat(contact.Username));

        const avatar = document.createElement("div");
        avatar.className = "avatar";

        const img = document.createElement("img");
        img.id = "avatar";
        img.src = "/static/images/avatar-white.png";
        avatar.appendChild(img);

        const name = document.createElement("span");
        name.textContent = contact.Username;

        const statusClass = contact.Online ? "online-dot" : "offline-dot";
        const status = document.createElement("div");
        status.className = statusClass;

        item.appendChild(avatar);
        item.appendChild(name);
        item.appendChild(status);
        friend.appendChild(item);

        div.appendChild(friend);
    });

    if (posts === 0) {
        div.innerHTML = `<div class="no-contacts"><p>No contacts available.</p></div>`;
    }

    if (typeof selectedUser !== "undefined" && selectedUser) {
        setChatStatusByUsername(selectedUser);
    }
}

function loadUser() {
    fetch("/api/me")
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                if (data.username != "") {
                    currentUsername = data.username;
                    let user = document.getElementById("username");
                    user.textContent = data.username;
                    let email = document.getElementById("email");
                    email.textContent = data.email;
                    document.getElementById("logout").classList.remove("hidden");
                    document.getElementById("auth-btns").classList.add("hidden");
                    document.getElementById("unauth-btns").classList.remove("hidden");
                    if (lastContacts) {
                        loadContacts(lastContacts);
                    }
                }
            } else {
                currentUsername = "";
                document.getElementById("logout").classList.add("hidden");
                document.getElementById("auth-btns").classList.remove("hidden");
            }
        })
        .catch(e => console.error(e));
}

const menu = document.getElementById("menu-dropdown");
const notifMenu = document.getElementById("notif-menu");

function toggleLogout(e) {
    e.stopPropagation();
    menu.classList.toggle("hidden");
    notifMenu.classList.add("hidden");
}

function toggleNotifs(e) {
    e.stopPropagation();
    const wasHidden = notifMenu.classList.contains("hidden");
    notifMenu.classList.toggle("hidden");
    menu.classList.add("hidden");

    // Mark as read on open but keep content visible.
    if (wasHidden && typeof markNotificationsRead === "function") {
        markNotificationsRead();
    }
}

document.addEventListener("click", () => {
    menu.classList.add("hidden");
    notifMenu.classList.add("hidden");
});
