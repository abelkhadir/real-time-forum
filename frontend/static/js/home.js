document.addEventListener("click", e => {
  if (e.target.dataset.action === "login") showLogin();
  if (e.target.dataset.action === "signup") showSignup();
});


function showLogin() {
    console.log("dddd")
}

function showSignup() {
    console.log("aaaa")
}