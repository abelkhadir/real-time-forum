function initApp() {
    // Initialize WebSocket first
    initWebSocket();

    // Load user info and contacts
    loadUser();

    // Setup post creation UI
    expandPostCreationArea();

    // Load initial posts
    getPosts(1);

    // Setup logout button
    setupLogoutButton();
}

function setupLogoutButton() {
    const logoutBtn = document.getElementById("logout");
    if (logoutBtn) {
        logoutBtn.addEventListener("click", handleLogout);
    }
}

function handleLogout() {
    fetch("/api/logout", {
        method: "POST",
        headers: { "Content-Type": "application/json" }
    })
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                showToast("green", "Logged out successfully");
                setTimeout(() => {
                    location.reload();
                }, 500);
            } else {
                showToast("red", "Logout failed");
            }
        })
        .catch(err => showToast("red", "Network error"));
}

// Start the app when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initApp);
} else {
    initApp();
}
