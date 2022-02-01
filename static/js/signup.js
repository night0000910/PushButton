$(document).ready(function() {
    var signupButton = document.getElementById("signup-btn");
    signupButton.addEventListener("click", signup);
})

function signup() {
    var username = document.getElementById("user-name").value;
    var password = document.getElementById("password").value;

    if (username && password) {
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
    var message = response.Message;

    if (message == "success") {
        window.location.href = "/succeed_in_signup"

    } else {

        var h6 = document.getElementById("message");
    
        if (h6 == null) {
            h6 = document.createElement("h6");
            h6.setAttribute("id", "message");
            var container = document.getElementById("container");
            var signupElement = document.getElementById("signup-div-element");
            container.insertBefore(h6, signupElement.nextSibling);

            h6.innerHTML = "既に同じ名前のユーザーが存在します"
        }
    }
}
