function getById(id) { 
    return document.getElementById(id); 
}

function fetchLogs() {

    fetch(`/api/log/pending`, { method: "GET" })
    .then((response) => response.json())
    .then((data) => {
        if (data.message == "not authenticated") {
            getById("message").style.display = 'block';
            getById("content").style.display = 'none';
        } else {
            console.log(data);
            getById("message").style.display = 'none';
            getById("content").style.display = 'block';
        }
    })
    .catch((err) => console.log(err))
}

function startLog(event) {

    event.preventDefault();

    // capture the startTime as soon as possible.
    const startTime = Math.floor(Date.now()/1000);
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(function (postion) {

            latitude = postion.coords.latitude;
            longitude = postion.coords.longitude;
            activity = getById("start-activity").value;
            notes = getById("start-notes").value;

            const endpoint = "/api/log/start/"
            fetch(host + endpoint, {
                method: "POST",
                body: JSON.stringify({
                    startTime,
                    latitude,
                    longitude,
                    activity,
                    notes
                })
            }).then((response) => response.json())
            .then((data) => {
                console.log(data);
            }).catch((err) => console.log(err));
        })
    }
}

function logout(event) {
    event.preventDefault();
    
    const endpoint = "/api/user/logout/"
    fetch(endpoint, {
        method: "POST",
    }).then((response) => response.json())
    .then((data) => {
        if (data.message == "ok") {
            localStorage.setItem("isAuthenticated", "false")
            redirectIfNeeded()
        }
    })
    .catch((err) => console.log(err))
}

function redirectIfNeeded() {
    
    item = localStorage.getItem("isAuthenticated")
    if (item == "true") {
        window.location.replace("./home.html")
    }
}

document.addEventListener('DOMContentLoaded', function() {
    fetchLogs();
    document.querySelector('#start-submit').addEventListener('click', startLog);
});