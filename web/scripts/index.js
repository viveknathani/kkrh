function getById(id) { 
    return document.getElementById(id); 
}

// if user is logged in, we move past the auth page
function redirectIfNeeded() {
    
    item = localStorage.getItem("isAuthenticated");
    if (item === "true") {
        window.location.replace("/web/home.html");
    } else {
        document.getElementsByTagName("body")[0].style.display = 'block';
    }
}

function login(event) {

    event.preventDefault();
    let email = getById("login-email").value;
    let password = getById("login-password").value;

    const endpoint = "/api/user/login/"
    fetch(endpoint, {
        method: "POST",
        body: JSON.stringify({
            email,
            password
        })
    }).then((response) => response.json())
    .then((data) => {
        if (data.message === "ok") {
            localStorage.setItem("isAuthenticated", "true")
            redirectIfNeeded();
        } else {
            msg = getById("login-message");
            msg.innerText = "oops! " + data.message;
            msg.style.color = 'red';
            
        }
    }).catch((err) => console.log(err));
}

function signup(event) {
    event.preventDefault();
    let name = getById("signup-name").value;
    let email = getById("signup-email").value;
    let password = getById("signup-password").value;

    const endpoint = "/api/user/signup/"
    fetch(endpoint, {
        method: "POST",
        body: JSON.stringify({
            name,
            email,
            password
        })
    }).then((response) => response.json())
    .then((data) => {
        if (data.message === "created") {
            msg = getById("signup-message");
            msg.innerText = data.message;
            msg.style.color = 'green';
        } else {
            msg = getById("signup-message");
            msg.innerText = "oops! " + data.message;
            msg.style.color = 'red';
            
        }
    }).catch((err) => console.log(err));
}

document.addEventListener('DOMContentLoaded', function() {
    redirectIfNeeded();
    document.querySelector('#signup-submit').addEventListener('click', signup);
    document.querySelector('#login-submit').addEventListener('click', login);
});