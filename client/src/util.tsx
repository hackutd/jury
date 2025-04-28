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

// Parse CSV data for preview (into string matrix)
// Sourced from https://stackoverflow.com/a/14991797
function parseCSV(str: string, delimiter: string): string[][] {
    const arr: string[][] = [];
    let quote = false;

    for (let row = 0, col = 0, c = 0; c < str.length; c++) {
        let cc = str[c],
            nc = str[c + 1]; // Current character, next character
        arr[row] = arr[row] || []; // Create a new row if necessary
        arr[row][col] = arr[row][col] || ''; // Create a new column (start with empty string) if necessary

        // If the current character is a quotation mark, and we're inside a
        // quoted field, and the next character is also a quotation mark,
        // add a quotation mark to the current column and skip the next character
        if (cc == '"' && quote && nc == '"') {
            arr[row][col] += cc;
            ++c;
            continue;
        }

        // If it's just one quotation mark, begin/end quoted field
        if (cc == '"') {
            quote = !quote;
            continue;
        }

        // If it's a comma and we're not in a quoted field, move on to the next column
        if (cc == ',' && !quote) {
            ++col;
            continue;
        }

        // If it's a newline (CRLF) and we're not in a quoted field, skip the next character
        // and move on to the next row and move to column 0 of that new row
        if (cc == '\r' && nc == '\n' && !quote) {
            ++row;
            col = 0;
            ++c;
            continue;
        }

        // If it's a newline (LF or CR) and we're not in a quoted field,
        // move on to the next row and move to column 0 of that new row
        if (cc == '\n' && !quote) {
            ++row;
            col = 0;
            continue;
        }
        if (cc == '\r' && !quote) {
            ++row;
            col = 0;
            continue;
        }

        // Otherwise, append the current character to the current column
        arr[row][col] += cc;
    }

    return arr;
}

// Convert a file to a string using FileReader
function fileToString(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = (event) => {
            if (event.target && event.target.result) {
                resolve(event.target.result as string);
            } else {
                reject(new Error('Failed to read file'));
            }
        };
        reader.onerror = (error) => {
            reject(error);
        };
        reader.readAsText(file);
    });
}

function pad(str: string, max: number) {
    if (str.length < max) {
        return str;
    }

    return str.substring(0, max) + '...';
}

export {
    timeSince,
    arrow,
    fixIfFloat,
    fixIfFloatDigits,
    errorAlert,
    showTopFive,
    parseCSV,
    fileToString,
    pad,
};
