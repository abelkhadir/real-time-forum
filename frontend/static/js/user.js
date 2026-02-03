function loadContacts(contacts, username) {
    if (!contacts || contacts.length === 0) {
        return;
    }

    let posts = 0

    const div = document.getElementById("friends-list");
    div.innerHTML = ""; // optional: clear old list

    contacts.forEach(contact => {
        if (username && contact.Username === username) return;
        posts++;
        const friend = document.createElement("div");
        friend.className = "friends-item";
        const statusClass = contact.Online ? "online-dot" : "offline-dot";
        friend.innerHTML = `
        <div class="friend-item" onclick="openChat('${contact.Username}')">
            <div class="avatar">
                <img id="avatar" src="/static/images/avatar-white.png">
            </div>
            <span>${contact.Username}</span>
            <div class="${statusClass}"></div>
        </div>
        `;

        div.appendChild(friend);
    });

    if (posts === 0) {
        div.innerHTML = `<div class="no-contacts"><p>No contacts available.</p></div>`;
    }

}

function loadUser() {
    fetch("/api/contacts")
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                if (data.username != "") {
                    let user = document.getElementById("username");
                    user.textContent = data.username;
                    let email = document.getElementById("email");
                    email.textContent = data.email;
                    document.getElementById("logout").classList.remove("hidden");
                    document.getElementById("auth-btns").classList.add("hidden");
                    document.getElementById("unauth-btns").classList.remove("hidden");
                }

                loadContacts(data.contacts, data.username);
            } else {
                document.getElementById("logout").classList.add("hidden");
                document.getElementById("auth-btns").classList.remove("hidden");
            }
        })
        .catch(e => console.error(e));
}


function toggleLogout() {
    document.getElementById("menu-dropdown").classList.toggle("hidden");
}