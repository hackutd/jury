window.onload = () => {
    // Only enable submit button if both checkboxes are checked
    document
        .getElementById("agreement-0")
        .addEventListener("change", checkChecked);
    document
        .getElementById("agreement-1")
        .addEventListener("change", checkChecked);
};

function checkChecked() {
    if (
        document.getElementById("agreement-0").checked &&
        document.getElementById("agreement-1").checked
    ) {
        document.getElementById("continue-button").disabled = false;
        return true;
    }
    document.getElementById("continue-button").disabled = true;
    return false;
}

async function login() {
    showLoginBlock();

    // Make sure both checkboxes are checked
    if (!checkChecked()) {
        alert("You must check both boxes!");
        hideLoginBlock();
        return;
    }

    // Make async call
    const res = await fetch("/judge/welcome", {
        method: "POST",
    });

    // Internal server error
    if (res.status !== 202) {
        const err = await res.text();
        console.error(err);
        alert(err);
        hideLoginBlock();
        return;
    }

    // Redirect
    window.location.href = window.origin + "/judge";

    hideLoginBlock();
}
