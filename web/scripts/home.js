import { getTotalTime } from "./common.js";

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
            getById("pending-message").style.display = 'none';
            localStorage.setItem("isAuthenticated", "false")
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

            const latitude = postion.coords.latitude;
            const longitude = postion.coords.longitude;
            const activity = getById("start-activity").value;
            const notes = getById("start-notes").value;

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

//if the user is not logged in, we need to go to the auth page
function redirectIfNeeded() {
    
    const item = localStorage.getItem("isAuthenticated");
    if (item === "false") {
        window.location.replace("/");
    }
}

// given startTime in UNIX format, get the time elapsed
// in hours:minutes:seconds format
function calculateTimeElapsed(startTime) {
    
    // factor in 1000 to convert to milliseconds as it works for the Date object
    startTime = startTime * 1000;
    const elapsed = (new Date().getTime() - startTime) / 1000;
    return getTotalTime(elapsed);
}

function applyStyleToTableOrRow(element) {
    element.style.border = '1px solid white';
    element.style.borderCollapse = 'collapse';
}

function fillTable(data) {

    let table = getById("pending");
    table.querySelectorAll('*').forEach(kid => kid.remove());

    applyStyleToTableOrRow(table)
    table.style.margin = 'auto';

    let store = new Map();
    for (let i = 0; i < data.length; ++i) {

        // create all elements
        let row = document.createElement('tr');
        let activityColumn = document.createElement('td');
        let elapsedColumn = document.createElement('td');
        let stopColumn = document.createElement('td');
        let stopButton = document.createElement('button');

        // bring in the text and append them
        activityColumn.innerText = data[i].activity;
        stopButton.innerText = 'stop';
        stopColumn.appendChild(stopButton);
        row.appendChild(activityColumn);
        row.appendChild(elapsedColumn);
        row.appendChild(stopColumn);
        table.appendChild(row);

        // sprinkle some CSS
        row.childNodes.forEach((kid, idx, parent) => {
            applyStyleToTableOrRow(kid)
            kid.style.padding = '10px';
        })

        // store element in Map
        store.set(elapsedColumn, 'go');

        // attach an interval that forms a closure with the
        // elapsedColumn and the Map and use these to alter
        // the text when Map shows that time is stopped
        let interval = setInterval((startTime) => {
            if (store.get(elapsedColumn) === 'stop') {
                clearInterval(interval);
            }
            elapsedColumn.innerText = calculateTimeElapsed(startTime);
        }, 1000, data[i].startTime);

        // setup a function to stop the time by forming a 
        // closure with the elapsedColumn
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
