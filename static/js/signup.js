$(document).ready(function() {
    var signupButton = document.getElementById("signup-btn");
    signupButton.addEventListener("click", );
})

function signup() {

    if (username && password) {
        var username = document.getElementById("user-name");
        var password = document.getElementById("password");
        var data = {username : username, password : password};
        var params = encodeHTMLForm(data);

        var request = new XMLHttpRequest();
        request.open("POST", "/signup");
        request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        request.onload = processAfterSignup;
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

function processAfterSignup() {
    var response = JSON.parse(this.response);
    console.log(response);
}
