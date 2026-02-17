// Toast Notification System
const container = document.getElementById('toast-container');

function showToast(color, message) {
    if (!container) return;

    if (container.children.length >= 4) {
        container.removeChild(container.firstChild);
    }

    const toast = document.createElement('div');
    toast.className = "toast";
    toast.style.backgroundColor = color;
    toast.innerText = message;

    container.appendChild(toast);

    setTimeout(() => toast.classList.add('hide'), 3000);
    setTimeout(() => toast.remove(), 3500);
}

function escapeHTML(value) {
    const str = String(value ?? "");
    return str
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/\"/g, "&quot;")
        .replace(/'/g, "&#39;");
}
