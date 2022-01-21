export function getTotalTime(elapsedSeconds) {

    let hours   = Math.floor((elapsedSeconds / 3600));
    let minutes = Math.floor((elapsedSeconds / 60) % 60);
    let seconds = Math.floor(elapsedSeconds % 60);

    // append a '0' if value is less than 10
    (hours < 10) ? hours = '0' + hours.toString() : hours.toString();
    (minutes < 10) ? minutes = '0' + minutes.toString() : minutes.toString();
    (seconds < 10) ? seconds = '0' + seconds.toString() : seconds.toString();
    return `${hours}:${minutes}:${seconds}`;
}