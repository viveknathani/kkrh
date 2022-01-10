const host = "https://kkrh.herokuapp.com"

function getById(id) { 
    return document.getElementById(id); 
}

function redirectIfNeeded() {
    
    item = localStorage.getItem("isAuthenticated")
    if (item == "true") {
        window.location.replace("./home.html")
    }
}

function login(event) {

    event.preventDefault();
    let email = getById("signup-email").value;
    let password = getById("signup-password").value;

    const endpoint = "/api/user/login/"
    fetch(host + endpoint, {
        method: "POST",
        body: JSON.stringify({
            email,
            password
        })
    }).then((response) => response.json())
    .then((data) => {
        if (data.message == "ok") {
            localStorage.setItem("isAuthenticated", "true")
            redirectIfNeeded()
        }
    }).catch((err) => console.log(err));
    fetch(host + "/api/log/pending", {method: "GET"})
}

function signup(event) {
    event.preventDefault();
    let name = getById("signup-name").value;
    let email = getById("signup-email").value;
    let password = getById("signup-password").value;

    const endpoint = "/api/user/signup/"
    fetch(host + endpoint, {
        method: "POST",
        body: JSON.stringify({
            name,
            email,
            password
        })
    }).then((response) => response.json())
    .then((data) => {
        console.log(data);
    }).catch((err) => console.log(err));
}

document.addEventListener('DOMContentLoaded', function() {
    redirectIfNeeded()
    document.querySelector('#signup-submit').addEventListener('click', signup);
    document.querySelector('#login-submit').addEventListener('click', login);
});