// Authentication Popups and Handlers

function openLoginPopup() {
    let div = document.getElementById("login-popup");
    div.innerHTML = `
    
    <div class="popup-overlay" onclick="closePopup()"></div>

    <div class="popup-content">

    <div class="auth-left">
        <h3>Welcome!</h3>
        <div class="side-svg">
            <svg width="80" height="80" viewBox="0 0 30 16">
                <path
                    d="m18.4 0-2.803 10.855L12.951 0H9.34L6.693 10.855 3.892 0H0l5.012 15.812h3.425l2.708-10.228 2.709 10.228h3.425L22.29 0h-3.892ZM24.77 13.365c0 1.506 1.12 2.635 2.615 2.635C28.879 16 30 14.87 30 13.365c0-1.506-1.12-2.636-2.615-2.636s-2.615 1.13-2.615 2.636Z">
                </path>
            </svg>
            <svg width="100" height="100" viewBox="0 0 163 163" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path
                    d="M81.09 162.18C125.875 162.18 162.18 125.875 162.18 81.09C162.18 36.3052 125.875 0 81.09 0C36.3053 0 0 36.3052 0 81.09C0 125.875 36.3053 162.18 81.09 162.18Z"
                    fill="#AAEEC4"></path>
                <path
                    d="M81.0906 147.041C117.514 147.041 147.041 117.514 147.041 81.0906C147.041 44.6674 117.514 15.1406 81.0906 15.1406C44.6675 15.1406 15.1406 44.6674 15.1406 81.0906C15.1406 117.514 44.6675 147.041 81.0906 147.041Z"
                    stroke="#502BD8" stroke-width="8" stroke-miterlimit="1.2"></path>
                <path
                    d="M103.819 78.3292C108.906 78.3292 113.029 74.1028 113.029 68.8892C113.029 63.6757 108.906 59.4492 103.819 59.4492C98.7329 59.4492 94.6094 63.6757 94.6094 68.8892C94.6094 74.1028 98.7329 78.3292 103.819 78.3292Z"
                    fill="#502BD8"></path>
                <path
                    d="M58.3662 78.3292C63.4528 78.3292 67.5762 74.1028 67.5762 68.8892C67.5762 63.6757 63.4528 59.4492 58.3662 59.4492C53.2797 59.4492 49.1562 63.6757 49.1562 68.8892C49.1562 74.1028 53.2797 78.3292 58.3662 78.3292Z"
                    fill="#502BD8"></path>
                <path
                    d="M48.8438 94.8906C52.3937 109.411 65.4838 120.181 81.0938 120.181C96.7037 120.181 109.794 109.411 113.344 94.8906"
                    stroke="#502BD8" stroke-width="8" stroke-miterlimit="1.2"></path>
            </svg>
        </div>
        <p>Not a member yet? <span onclick="openRegisterPopup()">Register now!</span></p>
    </div>

    <div class="auth-right">
        <h2>Log In</h2>
        <form id="login-form" class="login-form">
            <input type="identifier" placeholder="Email or Username" required>
            <input type="password" placeholder="Password" required>
            <button type="submit" class="btn btn-primary">Log in now</button>
            <div class="forgot"> 
                <p><span onclick="">Forgot your password?</span></p>
            </div>
        </form>
    </div>
    `;
    div.classList.remove("hidden");
    document.getElementById("login-form").addEventListener("submit", handleLogin);
}

function openRegisterPopup() {
    let div = document.getElementById("login-popup");
    div.innerHTML = `

    <div class="popup-overlay" onclick="closePopup()"></div>

    <div class="popup-content">

        <div class="auth-left">
            <h3>Welcome!</h3>
            <div class="side-svg">
                <svg width="80" height="80" viewBox="0 0 30 16">
                    <path
                        d="m18.4 0-2.803 10.855L12.951 0H9.34L6.693 10.855 3.892 0H0l5.012 15.812h3.425l2.708-10.228 2.709 10.228h3.425L22.29 0h-3.892ZM24.77 13.365c0 1.506 1.12 2.635 2.615 2.635C28.879 16 30 14.87 30 13.365c0-1.506-1.12-2.636-2.615-2.636s-2.615 1.13-2.615 2.636Z">
                    </path>
                </svg>
                <svg width="100" height="100" viewBox="0 0 163 163" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path
                        d="M81.09 162.18C125.875 162.18 162.18 125.875 162.18 81.09C162.18 36.3052 125.875 0 81.09 0C36.3052 0 0 36.3052 0 81.09C0 125.875 36.3052 162.18 81.09 162.18Z"
                        fill="#502BD8"></path>
                    <path
                        d="M81.0906 147.041C117.514 147.041 147.041 117.514 147.041 81.0906C147.041 44.6674 117.514 15.1406 81.0906 15.1406C44.6674 15.1406 15.1406 44.6674 15.1406 81.0906C15.1406 117.514 44.6674 147.041 81.0906 147.041Z"
                        stroke="#AAEEC4" stroke-width="8" stroke-miterlimit="1.2"></path>
                    <path
                        d="M103.812 78.3292C108.898 78.3292 113.022 74.1028 113.022 68.8892C113.022 63.6757 108.898 59.4492 103.812 59.4492C98.725 59.4492 94.6016 63.6757 94.6016 68.8892C94.6016 74.1028 98.725 78.3292 103.812 78.3292Z"
                        fill="#AAEEC4"></path>
                    <path
                        d="M58.3584 78.3292C63.445 78.3292 67.5684 74.1028 67.5684 68.8892C67.5684 63.6757 63.445 59.4492 58.3584 59.4492C53.2719 59.4492 49.1484 63.6757 49.1484 68.8892C49.1484 74.1028 53.2719 78.3292 58.3584 78.3292Z"
                        fill="#AAEEC4"></path>
                    <path
                        d="M48.8281 94.8906C52.3781 109.411 65.4681 120.181 81.0781 120.181C96.6881 120.181 109.778 109.411 113.328 94.8906"
                        stroke="#AAEEC4" stroke-width="8" stroke-miterlimit="1.2"></path>
                </svg>
            </div>
            <p>Are you a member? <span onclick="openLoginPopup()">Log in now!</span></p>
        </div>

        <div class="auth-right">
            <h2>Sign up</h2>
            <form id="register-form" class="login-form">
                <input type="text" id="username" placeholder="Username" required>
                <input type="email" placeholder="Email" required>
                <input type="password" placeholder="Password" required>
                <input type="password" placeholder="Confirm Password" required>
                <button type="submit" class="btn btn-primary">Sign up now</button>
            </form>
        </div>

    </div>
    `;
    div.classList.remove("hidden");
    document.getElementById("register-form").addEventListener("submit", handleRegister);
}

function closePopup() {
    document.getElementById("login-popup").classList.add("hidden");
}

function handleLogin(e) {
    e.preventDefault();
    const form = e.target;
    const identifier = form.querySelector("input[type='identifier']").value;
    const password = form.querySelector("input[type='password']").value;

    fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ identifier, password })
    })
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                showToast("green", "Logged in successfully");
                closePopup();
                setTimeout(() => {
                    location.reload();
                }, 500);

                document.getElementById("logout").classList.remove("hidden");
                document.getElementById("auth-btns").classList.add("hidden");
            } else {
                showToast("red", data.error || "Login failed");
            }
        })
        .catch(err => showToast("red", "Network error"));
}

function handleRegister(e) {
    e.preventDefault();
    const form = e.target;
    const username = form.querySelector("input[type='text']").value;
    const email = form.querySelector("input[type='email']").value;
    const password = form.querySelectorAll("input[type='password']")[0].value;
    const confirmPassword = form.querySelectorAll("input[type='password']")[1].value;

    if (password !== confirmPassword) {
        showToast("red", "Passwords do not match");
        return;
    }

    fetch("/api/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, email, password })
    })
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                showToast("green", "Account created successfully");
                openLoginPopup();
            } else {
                showToast("red", data.error || "Registration failed");
            }
        })
        .catch(err => showToast("red", "Network error"));
}
