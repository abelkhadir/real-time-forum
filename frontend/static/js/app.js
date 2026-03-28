// initApp loads the authenticated app state and realtime features.
async function initApp() {
    // Setup logout button
    setupLogoutButton();

    const authenticated = await loadUser();
    if (!authenticated) {
        return;
    }

    // Initialize WebSocket first
    initWebSocket();

    // Load user info and contacts
    fetchNotifications();

    // Setup post creation UI
    expandPostCreationArea();
    if (typeof bindPostsPagination === "function") {
        bindPostsPagination();
    }

    // Load initial posts
    getPosts(1);
}

// setupLogoutButton wires the logout button click handler.
function setupLogoutButton() {
    const logoutBtn = document.getElementById("logout");
    if (logoutBtn) {
        logoutBtn.addEventListener("click", handleLogout);
    }
}

// handleLogout ends the current session and reloads the page.
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
