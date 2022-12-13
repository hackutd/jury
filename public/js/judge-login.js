let loginLock = false;

window.onload = () => {
    // Clear the judge login cookie if invalid
    deleteCookie("token");

    // Add event listener to keydown for login input
    document
        .getElementById("login-code-input")
        .addEventListener("keydown", (key) => {
            if (key.key == "Enter") {
                login();
            }
            hideError();
        });

    // Only enable submit button if 6 digits entered
    document
        .getElementById("login-code-input")
        .addEventListener("input", (v) => {
            if (document.getElementById("login-code-input").value.length < 6) {
                document.getElementById("submit-button").disabled = true;
            } else {
                document.getElementById("submit-button").disabled = false;
            }
        });
};

async function login() {
    // Prevent spamming of login while waiting for fetch
    if (loginLock) return;
    showLoginBlock();

    // Get code from input
    const code = document.getElementById("login-code-input").value;

    // Check for length of code
    if (code.length < 6) {
        showError();
        hideLoginBlock();
        return;
    }

    // Make async call to check code
    const res = await fetch("/judge/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ code }),
    });

    // Invalid code
    if (res.status === 400) {
        showError();
        hideLoginBlock();
        return;
    }

    // Internal server error
    if (res.status !== 200) {
        const err = await res.text();
        console.error(err);
        alert(err);
        showError();
        hideLoginBlock();
        return;
    }

    // Correct code; save token as cookie
    const token = await res.text();
    setCookie("token", token, 60 * 60 * 24);

    // Redirect
    window.location.href = window.origin + "/judge";

    hideLoginBlock();
}

function showError() {
    document.getElementById("login-code-input").classList.add("error");
    document.getElementById("login-code-input-label").classList.add("error");
    document.getElementById("login-code-input-label").textContent =
        "Invalid judging code";
}

function hideError() {
    document.getElementById("login-code-input").classList.remove("error");
    document.getElementById("login-code-input-label").classList.remove("error");
    document.getElementById("login-code-input-label").textContent =
        "Enter your judging code";
}
