$(document).ready(function() {
    var loginButton = document.getElementById("login-btn");
    loginButton.addEventListener("click", login);
})

function login() {
    var username = document.getElementById("user-name").value;
    var password = document.getElementById("password").value;

    if (username && password) {
        var data = {username : username, password : password}
        var params = encodeHTMLForm(data);

        var request = new XMLHttpRequest();
        request.open("POST", "/login");
        request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        request.onload = processAfterLogin;
        request.send(params);

    }
}

function encodeHTMLForm(data) {
    var params = [];

    for (var name in data) {
        var value = data[name];
        var param = encodeURIComponent(name) + "=" + encodeURIComponent(value);
        params.push(param);
    }

    return params.join("&").replace(/%20/g, "+");
}

function processAfterLogin() {
    var response = JSON.parse(this.response);
    var message = response.Message;

    if (message == "success") {
        window.location.href = "/push"

    } else if (message == "failed") {
        var h6 = document.getElementById("message");
    
        if (h6 == null) {
            h6 = document.createElement("h6");
            h6.setAttribute("id", "message");
            var container = document.getElementById("container");
            var loginElement = document.getElementById("login-div-element");
            container.insertBefore(h6, loginElement.nextSibling);

            h6.innerHTML = "ログインに失敗しました";
        }

    } else if (message == "authenticated") {
        window.location.href = "/push"
    }
}