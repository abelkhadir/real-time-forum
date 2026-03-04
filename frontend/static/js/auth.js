// Authentication popups and handlers

function notify(color, message) {
    if (typeof showToast === "function") {
        showToast(color, message);
        return;
    }

    if (color === "red") {
        console.error(message);
        return;
    }

    console.log(message);
}

function openLoginPopup() {
    const loginPopup = document.getElementById("login-popup");
    const registerPopup = document.getElementById("register-popup");

    if (registerPopup) {
        registerPopup.classList.add("hidden");
    }

    if (loginPopup) {
        loginPopup.classList.remove("hidden");
    }
}

function openRegisterPopup() {
    const loginPopup = document.getElementById("login-popup");
    const registerPopup = document.getElementById("register-popup");

    if (loginPopup) {
        loginPopup.classList.add("hidden");
    }

    if (registerPopup) {
        registerPopup.classList.remove("hidden");
    }
}

function closePopup() {
    const loginPopup = document.getElementById("login-popup");
    const registerPopup = document.getElementById("register-popup");

    if (loginPopup) {
        loginPopup.classList.add("hidden");
    }

    if (registerPopup) {
        registerPopup.classList.add("hidden");
    }
}

async function handleLogin(e) {
    e.preventDefault();

    const identifierInput = document.getElementById("login-identifier");
    const passwordInput = document.getElementById("login-password");

    const identifier = (identifierInput?.value || "").trim();
    const password = passwordInput?.value || "";

    if (!identifier || !password) {
        notify("red", "Please fill in both fields");
        return;
    }

    try {
        const res = await fetch("/api/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ identifier, password }),
        });

        const data = await res.json().catch(() => ({}));
        if (data.success) {
            notify("green", "Logged in successfully");

            closePopup();
            window.location.assign("/");
            return;
        }

        notify("red", data.error || "Login failed");
    } catch (err) {
        notify("red", "Network error");
    }
}

async function handleRegister(e) {
    e.preventDefault();

    const usernameInput = document.getElementById("register-username");
    const ageInput = document.getElementById("register-age");
    const genderInput = document.getElementById("register-gender");
    const firstNameInput = document.getElementById("register-first-name");
    const lastNameInput = document.getElementById("register-last-name");
    const emailInput = document.getElementById("register-email");
    const passwordInput = document.getElementById("register-password");
    const confirmInput = document.getElementById("register-confirm-password");

    const username = (usernameInput?.value || "").trim();
    const age = Number(ageInput?.value || 0);
    const gender = (genderInput?.value || "").trim();
    const first_name = (firstNameInput?.value || "").trim();
    const last_name = (lastNameInput?.value || "").trim();
    const email = (emailInput?.value || "").trim();
    const password = passwordInput?.value || "";
    const confirmPassword = confirmInput?.value || "";

    if (!username || !age || !gender || !first_name || !last_name || !email || !password || !confirmPassword) {
        notify("red", "Please complete all fields");
        return;
    }

    if (password !== confirmPassword) {
        notify("red", "Passwords do not match");
        return;
    }

    try {
        const res = await fetch("/api/register", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                username,
                age,
                gender,
                first_name,
                last_name,
                email,
                password
            }),
        });

        const data = await res.json().catch(() => ({}));
        if (data.success) {
            notify("green", "Account created successfully");
            const registerForm = document.getElementById("register-form");
            if (registerForm) {
                registerForm.reset();
            }
            openLoginPopup();
            return;
        }

        notify("red", data.error || "Registration failed");
    } catch (err) {
        notify("red", "Network error");
    }
}

function bindAuthForms() {
    const loginForm = document.getElementById("login-form");
    if (loginForm && loginForm.dataset.bound !== "true") {
        loginForm.addEventListener("submit", handleLogin);
        loginForm.dataset.bound = "true";
    }

    const registerForm = document.getElementById("register-form");
    if (registerForm && registerForm.dataset.bound !== "true") {
        registerForm.addEventListener("submit", handleRegister);
        registerForm.dataset.bound = "true";
    }
}

if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", bindAuthForms);
} else {
    bindAuthForms();
}
