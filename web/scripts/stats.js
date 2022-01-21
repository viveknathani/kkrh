import { getTotalTime } from "./common.js";

function getById(id) { 
    return document.getElementById(id); 
}

// initiate a pie chart with some colours
function prepareAndSendForDrawing(store, totalSeconds) {

    let arr = [];
    let colorSet = ["#FF6633", "#FFB399", "#FF33FF", "#FFFF99", "#00B3E6", "#E6B333",
    "#3366E6", "#999966", "#809980", "#E6FF80", "#1AFF33", "#999933", "#FF3380",
    "#CCCC00", "#66E64D", "#4D80CC", "#FF4D4D", "#99E6E6", "#6666FF" ];
    
    let sum = 0;
    store.forEach((value, key) => {
        let temp = ((value * 100.0)) / Number(totalSeconds);
        sum += temp;
        arr.push({ activity: key, timeSpent: temp });
    });
    store.set('unrecorded', (100.0 - sum) * totalSeconds / 100.0)
    arr.push({activity: 'unrecorded', timeSpent: 100 - sum });

    let labels = [];
    let colors = [];
    let data = [];

    for (let i = 0; i < arr.length; ++i) {
        labels.push(`${arr[i].activity}(${getTotalTime(store.get(arr[i].activity))})`);
        data.push(arr[i].timeSpent);
        colors.push(colorSet[i % colorSet.length]);
    }

    makePieChart(labels, data, colors);
}

// make the API call and send data for drawing
function getLogsInRange(startTime, endTime) {

    fetch(`/api/stats/range?startTime=${startTime}&endTime=${endTime}`, {
        method: 'GET',
    }).then((response) => response.json())
    .then((data) => {

        const totalSeconds = endTime - startTime;
        let store = new Map();
        for (let i = 0; i < data.length; ++i) {
            const prev = store.get(data[i].activity);
            let value = data[i].endTime - data[i].startTime;
            if (prev != undefined) {
                value += prev;
            }
            store.set(data[i].activity, value);
        }

        prepareAndSendForDrawing(store, totalSeconds);

    }).catch((err) => console.log(err));
}

// runs when "get logs" button is pressed
function getLogsInRangeHandler(event) {

    event.preventDefault();
    let startDate = getById("stats-start").value;
    let endDate = getById("stats-end").value;
    let startTime = Date.parse(startDate)/1000;
    let endTime = Date.parse(endDate)/1000;
    getLogsInRange(startTime, endTime);
}

function makePieChart(labels, store, colors) {

    let chartStatus = Chart.getChart("board"); 
    if (chartStatus != undefined) {
        chartStatus.destroy();
    }

    const data = {
        labels: labels,
        datasets: [{
            label: 'pie chart',
            data: store,
            backgroundColor: colors,
            hoverOffset: 4
        }],
    };
    const config = {
        type: 'pie',
        aspectRatio: 1,
        data: data,
        options: {
            responsive: true,
            radius: 150,
            legend: {
                display: true,
                position: "bottom",
                labels: {
                  fontColor: "white",
                  fontSize: 16
                }    
            },
        },
    };
    const board = getById("board");
    const ctx = board.getContext("2d");
    const myPie = new Chart(ctx, config);
}

document.addEventListener('DOMContentLoaded', function() {
    document.querySelector('#stats-submit').addEventListener('click', getLogsInRangeHandler);
});
