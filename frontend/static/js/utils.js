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
