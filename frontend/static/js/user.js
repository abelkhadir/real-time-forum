let contactsByName = {};
let currentUsername = "";
let lastContacts = null;

function setAuthenticatedState(isAuthenticated) {
    const appShell = document.getElementById("app-shell");
    const authButtons = document.getElementById("auth-btns");
    const profileButtons = document.getElementById("unauth-btns");

    if (appShell) {
        appShell.classList.toggle("hidden", !isAuthenticated);
    }

    if (authButtons) {
        authButtons.classList.toggle("hidden", isAuthenticated);
    }

    if (profileButtons) {
        profileButtons.classList.toggle("hidden", !isAuthenticated);
    }

    if (isAuthenticated) {
        closePopup();
    } else {
        openLoginPopup();
    }
}

function setChatStatusByUsername(username) {
    const statusEl = document.getElementById("chat-status");
    if (!statusEl) return;
    const contact = contactsByName[username];
    const online = contact ? contact.Online : false;
    statusEl.textContent = online ? "Online" : "Offline";
    statusEl.classList.toggle("online", online);
    statusEl.classList.toggle("offline", !online);
}

function loadContacts(contacts, wsUsername = "") {
    if (!contacts || contacts.length === 0) {
        return;
    }

    lastContacts = contacts;
    const selfUsername = (currentUsername || wsUsername || "").trim();

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
    return fetch("/api/me")
        .then(async (res) => {
            if (!res.ok) {
                return { success: false };
            }
            return res.json();
        })
        .then(data => {
            if (data.success) {
                if (data.username != "") {
                    currentUsername = data.username;
                    let user = document.getElementById("username");
                    user.textContent = data.username;
                    let email = document.getElementById("email");
                    email.textContent = data.email;
                    document.getElementById("logout").classList.remove("hidden");
                    setAuthenticatedState(true);
                    if (lastContacts) {
                        loadContacts(lastContacts);
                    }
                    return true;
                }
            } else {
                currentUsername = "";
                document.getElementById("logout").classList.add("hidden");
                setAuthenticatedState(false);
                return false;
            }
            setAuthenticatedState(false);
            return false;
        })
        .catch(e => {
            console.error(e);
            currentUsername = "";
            setAuthenticatedState(false);
            return false;
        });
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
