function getCookies() {
    const cookieList = document.cookie.split(";");
    const splitCookies = cookieList.reduce((acc, curr) => {
        const cs = curr.split("=");
        acc[cs[0]] = cs[1];
        return acc;
    }, {});
    return splitCookies;
}

function getCookie(key) {
    return getCookies()[key];
}

// Tokens should get 60*60*24 = 1 day
function setCookie(key, value, maxAge = null) {
    const maxAgeString = maxAge ? `maxage=${maxAge}; ` : "";
    document.cookie = `${key}=${value}; ${maxAgeString}path=/; samesite=strict; secure;`;
}

function deleteCookie(key) {
    document.cookie = `${key}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; samesite=strict; secure;`;
}

function showLoginBlock() {
    document.getElementsByClassName("login-block")[0].style.display = "flex";
    loginLock = true;
}

function hideLoginBlock() {
    document.getElementsByClassName("login-block")[0].style.display = "none";
    loginLock = false;
}

function logoutJudge() {
    deleteCookie("token");
    window.location.href = window.origin;
}

function logoutAdmin() {
    deleteCookie("admin-pass");
    window.location.href = window.origin;
}

function getQueryParams() {
    const params = new Proxy(new URLSearchParams(window.location.search), {
        get: (searchParams, prop) => searchParams.get(prop),
    });
    return params;
}