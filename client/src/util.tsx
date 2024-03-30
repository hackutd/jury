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
    if (!n) {
        return '0';
    }
    if (Math.round(n) === n) {
        return n.toString();
    }
    return n.toFixed(d);
}

// Show the top five elements of an array, with an ellipsis if there are more
function showTopFive<T>(arr: T[]): string {
    if (arr.length === 0) {
        return '[none]';
    }
    return arr.slice(0, 5).join(', ') + (arr.length > 5 ? '...' : '');
}

function errorAlert<T>(res: FetchResponse<T>) {
    const err = `Error sending request to server (Status ${res.status}): ${res.error}`;
    alert(err);
    console.error(err);
}

export { timeSince, arrow, fixIfFloat, fixIfFloatDigits, errorAlert, showTopFive };
