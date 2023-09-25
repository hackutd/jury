// Convert millisecond time to "x secs/mins/hours ago"
const timeSince = (date: number) => {
    // eslint-disable-next-line eqeqeq
    if (date == 0) {
        return 'never';
    }
    const seconds = Math.floor((new Date().getTime() - date) / 1000);
    if (seconds < 0) {
        return Math.abs(seconds) + ' seconds in the future?!?!';
    }
    let interval = seconds / 31536000;
    if (interval > 1) {
        return Math.floor(interval) + ' years ago';
    }
    interval = seconds / 2592000;
    if (interval > 1) {
        return Math.floor(interval) + ' months ago';
    }
    interval = seconds / 86400;
    if (interval > 1) {
        return Math.floor(interval) + ' days ago';
    }
    interval = seconds / 3600;
    if (interval > 1) {
        return Math.floor(interval) + ' hours ago';
    }
    interval = seconds / 60;
    if (interval > 1) {
        return Math.floor(interval) + ' minutes ago';
    }
    return 'just now';
};

const arrow = (asc: boolean) => (asc ? '▲' : '▼');

function fixIfFloat(n: number): string {
    return fixIfFloatDigits(n, 3);
}

function fixIfFloatDigits(n: number, d: number): string {
    if (Math.round(n) === n) {
        return n.toString();
    }
    return n.toFixed(d);
}

export { timeSince, arrow, fixIfFloat, fixIfFloatDigits };
