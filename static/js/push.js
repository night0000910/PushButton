$(document).ready(function() {
    var pushButton = document.getElementById("push-btn");
    pushButton.addEventListener("click", earnMoney);
    var resetButton = document.getElementById("reset-btn");
    resetButton.addEventListener("click", reset);
    var investButton = document.getElementById("invest-btn");
    investButton.addEventListener("click", invest);
})

function getCsrfToken() {
    var csrfToken = document.getElementById("csrf-token").value;
    return csrfToken
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

function earnMoney() {
    var csrfToken = getCsrfToken();
    var data = {csrfToken : csrfToken};
    params = encodeHTMLForm(data)

    var request = new XMLHttpRequest();
    request.open("POST", "/earn_money");
    request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    request.onload = processAfterEarnMoney;
    request.send(params);
}

function processAfterEarnMoney() {
    var response = JSON.parse(this.response);
    var message = response.Message;
    var money = response.Money;
    var moneyElement = document.getElementById("money");
    moneyElement.innerHTML = money;

    if (message == "not authenticated") {
        window.location.href = "/homepage";
    }
}

function invest() {
    var csrfToken = getCsrfToken();
    var data = {csrfToken : csrfToken};
    params = encodeHTMLForm(data)

    var request = new XMLHttpRequest();
    request.open("POST", "/invest");
    request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    request.onload = processAfterInvest;
    request.send(params);
}

function processAfterInvest() {
    var response = JSON.parse(this.response);
    var message = response.Message;
    var money = response.Money;
    var thunderFruits = response.ThunderFruits;
    var profit = response.Profit;

    var moneyElement = document.getElementById("money");
    var thunderFruitsElement = document.getElementById("thunder-fruits");
    var pushButton = document.getElementById("push-btn");
    var h5 = document.getElementById("oak-h5-element");
    var container = document.getElementById("container");

    var messageList = message.split("-");
    message = messageList[0];

    if (message == "invested") {
        var investMessage = messageList[1];

        moneyElement.innerHTML = money;
        thunderFruitsElement.innerHTML = thunderFruits;

        var h6 = document.getElementById("message");

        if (h6 == null) {
            var h6 = document.createElement("h6");
            h6.setAttribute("id", "message");
            container.insertBefore(h6, h5);
        }

        if (investMessage == "success") {
            h6.setAttribute("style", "color: #DC143C;");
            h6.innerHTML = profit + "$?????????????????????!!"
        } else {
            h6.setAttribute("style", "color: #333;");
            h6.innerHTML = "?????????0$?????????..."
        }

    } else if (message == "not invested") {
        var h6 = document.getElementById("message");

        if (h6 == null) {
            var h6 = document.createElement("h6");
            h6.setAttribute("id", "message");
            container.insertBefore(h6, h5);
        }

        h6.setAttribute("style", "color: #333;");
        h6.innerHTML = "????????????????????????????????????"

    } else if (message == "not authenticated") {
        window.location.href = "/homepage"
    }
}

function reset() {
    var csrfToken = getCsrfToken();
    var data = {csrfToken : csrfToken};
    params = encodeHTMLForm(data)

    var request = new XMLHttpRequest();
    request.open("POST", "/reset");
    request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    request.onload = processAfterReset;
    request.send(params);
}

function processAfterReset() {
    var response = JSON.parse(this.response);
    var message = response.Message;
    var money = response.Money;
    var oakFruits = response.OakFruits;
    var thunderFruits = response.ThunderFruits;

    var moneyElement = document.getElementById("money");
    moneyElement.innerHTML = money;
    var oakFruitsElement = document.getElementById("oak-fruits");
    oakFruitsElement.innerHTML = oakFruits;
    var thunderFruitsElement = document.getElementById("thunder-fruits");
    thunderFruitsElement.innerHTML = thunderFruits;

    if (message == "not authenticated") {
        window.location.href = "/homepage";
    }
}