let loginLock = false;

window.onload = () => {
    // Clear the admin login cookie if invalid
    deleteCookie('admin-pass');

    // Add listener for password input box
    document.getElementById('login-code-input').addEventListener('keydown', (key) => {
        if (key.key == 'Enter') {
            login();
        }
        hideError();
    });
};

async function login() {
    // Prevent spamming of login while waiting for fetch
    if (loginLock) return;
    loginLock = true;

    // Get code from input
    const password = document.getElementById('login-code-input').value;

    // Make async call to check code
    const res = await fetch('/admin/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ password }),
    });

    // Invalid code
    if (res.status === 400) {
        showError();
        loginLock = false;
        return;
    }

    // Internal server error
    if (res.status !== 200) {
        console.error(await res.text());
        showError();
        loginLock = false;
        return;
    }

    // Correct code; save password as cookie
    setCookie('admin-pass', password, 60 * 60 * 24);

    // Redirect
    window.location.href = window.origin + '/admin';

    loginLock = false;
}

function showError() {
    document.getElementById('login-code-input').classList.add('error');
    document.getElementById('login-code-input-label').classList.add('error');
    document.getElementById('login-code-input-label').textContent = 'Password is invalid';
}

function hideError() {
    document.getElementById('login-code-input').classList.remove('error');
    document.getElementById('login-code-input-label').classList.remove('error');
    document.getElementById('login-code-input-label').textContent = 'Admin Password';
}
