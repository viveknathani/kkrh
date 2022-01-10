const host = "https://kkrh.herokuapp.com"

function getById(id) { 
    return document.getElementById(id); 
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

document.addEventListener('DOMContentLoaded', function() {
    document.querySelector('#start-submit').addEventListener('click', startLog);
});