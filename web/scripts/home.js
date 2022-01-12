function getById(id) { 
    return document.getElementById(id); 
}

function fetchLogs() {

    getById("pending-message").innerText = 'fetching logs...';
    const endoint = "/api/log/pending"
    fetch(endoint, { method: "GET" })
    .then((response) => response.json())
    .then((data) => {
        if (data.message == "not authenticated") {
            getById("message").style.display = 'block';
            getById("content").style.display = 'none';
        } else {
            if (data.length === 0) {
                getById("pending-message").innerText = 'you have no pending logs';
            } else {
                getById("pending-message").innerText = 'here are your pending logs';
            }
            fillTable(data);
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
            fetch(endpoint, {
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
                fetchLogs();
            }).catch((err) => console.log(err));
        })
    }
}

function endLog(logId, endTime) {

    const endpoint = "/api/log/stop/"
    fetch(endpoint, {
        method: "PUT",
        body: JSON.stringify({
            logId,
            endTime
        })
    }).then((response) => response.json())
    .then((data) => {
        console.log(data);
        fetchLogs();
    }).catch((err) => console.log(err));

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
    
    const item = localStorage.getItem("isAuthenticated");
    if (item === "false") {
        window.location.replace("/");
    }
}

function calculateTimeElapsed(startTime) {
    
    startTime = startTime * 1000;
    const elapsed = (new Date().getTime() - startTime) / 1000;
    let hours   = Math.floor((elapsed / 3600) % 24);
    let minutes = Math.floor((elapsed / 60) % 60);
    let seconds = Math.floor(elapsed % 60);

    (hours < 10) ? hours = '0' + hours.toString() : hours.toString();
    (minutes < 10) ? minutes = '0' + minutes.toString() : minutes.toString();
    (seconds < 10) ? seconds = '0' + seconds.toString() : seconds.toString();
    return `${hours}:${minutes}:${seconds}`;
}

function fillTable(data) {

    let table = getById("pending");
    table.querySelectorAll('*').forEach(kid => kid.remove());

    table.style.border = '1px solid white';
    table.style.borderCollapse = 'collapse';
    table.style.margin = 'auto';

    let store = new Map();
    for (let i = 0; i < data.length; ++i) {

        let row = document.createElement('tr');
        let activityColumn = document.createElement('td');
        let elapsedColumn = document.createElement('td');
        let stopColumn = document.createElement('td');
        let stopButton = document.createElement('button');

        activityColumn.innerText = data[i].activity;
        stopButton.innerText = 'stop';
        stopColumn.appendChild(stopButton);
        row.appendChild(activityColumn);
        row.appendChild(elapsedColumn);
        row.appendChild(stopColumn);
        table.appendChild(row);

        row.childNodes.forEach((kid, idx, parent) => {
            kid.style.border = '1px solid white';
            kid.style.borderCollapse = 'collapse';
            kid.style.padding = '10px';
        })

        store.set(elapsedColumn, 'go');
        let interval = setInterval((startTime) => {
            if (store.get(elapsedColumn) === 'stop') {
                clearInterval(interval);
            }
            elapsedColumn.innerText = calculateTimeElapsed(startTime);
        }, 1000, data[i].startTime);

        stopButton.onclick = function(event) {
            event.preventDefault();
            store.set(elapsedColumn, 'stop');
            endLog(data[i].id, Math.floor(Date.now()/1000));
        }
    }
}

document.addEventListener('DOMContentLoaded', function() {
    fetchLogs();
    document.querySelector('#start-submit').addEventListener('click', startLog);
    document.querySelector('#logout').addEventListener('click', logout);
});
