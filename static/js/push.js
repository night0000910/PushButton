$(document).ready(function() {
    var pushButton = document.getElementById("push-btn");
    pushButton.addEventListener("click", earnMoney);
    var resetButton = document.getElementById("reset-btn");
    resetButton.addEventListener("click", reset);
    var investButton = document.getElementById("invest-btn");
    investButton.addEventListener("click", invest);
})

function earnMoney() {
    var request = new XMLHttpRequest();
    request.open("POST", "/earn_money");
    request.onload = processAfterEarnMoney;
    request.send();
}

function processAfterEarnMoney() {
    var response = JSON.parse(this.response);
    var money = response.Money;
    var moneyElement = document.getElementById("money");
    moneyElement.innerHTML = money;
}

function invest() {
    var request = new XMLHttpRequest();
    request.open("POST", "/invest");
    request.onload = processAfterInvest;
    request.send();
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
            h6.innerHTML = profit + "$手に入れました!!"
        } else {
            h6.setAttribute("style", "color: #333;");
            h6.innerHTML = "利益は0$でした..."
        }

    } else {
        var h6 = document.getElementById("message");

        if (h6 == null) {
            var h6 = document.createElement("h6");
            h6.setAttribute("id", "message");
            container.insertBefore(h6, h5);
        }

        h6.setAttribute("style", "color: #333;");
        h6.innerHTML = "サンダーの実がありません"
    }
}

function reset() {
    var request = new XMLHttpRequest();
    request.open("POST", "/reset");
    request.onload = processAfterReset;
    request.send();
}

function processAfterReset() {
    var response = JSON.parse(this.response);
    var money = response.Money;
    var oakFruits = response.OakFruits;
    var thunderFruits = response.ThunderFruits;

    var moneyElement = document.getElementById("money");
    moneyElement.innerHTML = money;
    var oakFruitsElement = document.getElementById("oak-fruits");
    oakFruitsElement.innerHTML = oakFruits;
    var thunderFruitsElement = document.getElementById("thunder-fruits");
    thunderFruitsElement.innerHTML = thunderFruits;
}