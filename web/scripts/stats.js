class PieChart {

    constructor(canvas, data, colors) {
        this.canvas = canvas;
        this.data = data;
        this.colors = colors;
    }

    drawPieSlice(ctx, startAngle, endAngle, color, text) {
        let centerX = 150;
        let centerY = 150;
        let radius = 150;
        ctx.fillStyle = color;
        ctx.beginPath();
        ctx.moveTo(centerX,centerY);
        ctx.arc(centerX, centerY, radius, startAngle, endAngle);
        ctx.closePath();
        ctx.fill();

        let labelText = text;
        let labelX = (this.canvas.width / 2) + (radius / 2) * Math.cos(startAngle + (endAngle - startAngle) / 2);
        let labelY = (this.canvas.height / 2) + (radius / 2) * Math.sin(startAngle + (endAngle - startAngle) / 2);
        ctx.fillStyle = "white";
        ctx.font = "bold 10px Arial";
        ctx.fillText(labelText, labelX, labelY);
    }

    draw() {

        const ctx = this.canvas.getContext("2d");
        let startAngle = 0;
        for (let i = 0; i < this.data.length; ++i) {
            const angle = 2 * Math.PI * this.data[i].timeSpent / 100.0;
            this.drawPieSlice(ctx, startAngle, startAngle + angle, this.colors[i % this.colors.length], i.toString());
            startAngle += angle;
            addParagraph(`${i.toString()}:${this.data[i].activity}`);
        }
    }
}

function addParagraph(text) {
    let body = document.getElementsByTagName('body')[0];
    let p = document.createElement('p');
    p.innerText = text;
    body.appendChild(p);
} 

function getById(id) { 
    return document.getElementById(id); 
}

function prepareAndSendForDrawing(store, totalSeconds) {

    let arr = [];
    let colors = ["#FF6633", "#FFB399", "#FF33FF", "#FFFF99", "#00B3E6", "#E6B333",
    "#3366E6", "#999966", "#809980", "#E6FF80", "#1AFF33", "#999933", "#FF3380",
    "#CCCC00", "#66E64D", "#4D80CC", "#FF4D4D", "#99E6E6", "#6666FF" ];
    
    let sum = 0;
    store.forEach((value, key) => {
        let temp = ((value * 100.0)) / Number(totalSeconds);
        sum += temp;
        arr.push({ activity: key, timeSpent: 40 });
    });
    arr.push({activity: 'unrecorded', timeSpent: 20 });

    let board = getById("board");
    board.width = 300;
    board.height = 300;
    let pie = new PieChart(board, arr, colors)
    pie.draw();
}

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

function getLogsInRangeHandler(event) {

    event.preventDefault();
    let startDate = getById("stats-start").value;
    let endDate = getById("stats-end").value;
    let startTime = Date.parse(startDate)/1000;
    let endTime = Date.parse(endDate)/1000;
    getLogsInRange(startTime, endTime);
}

document.addEventListener('DOMContentLoaded', function() {
    document.querySelector('#stats-submit').addEventListener('click', getLogsInRangeHandler);
});
