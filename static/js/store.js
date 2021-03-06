$(document).ready(function() {
    displayPriceOfFruits();

    var oakFruitsButton = document.getElementById("oak-fruits-btn");
    oakFruitsButton.addEventListener("click", buyOakFruits);
    var thunderFruitsButton = document.getElementById("thunder-fruits-btn");
    thunderFruitsButton.addEventListener("click", buyThunderFruits);
})

function displayPriceOfFruits() {
    var request = new XMLHttpRequest();
    request.open("GET", "/users_information");
    request.onload = insertPriceIntoElements;
    request.send();
}

function insertPriceIntoElements() {
    var response = JSON.parse(this.response);
    var priceOfOakFruits = response.PriceOfOakFruits;
    var priceOfThunderFruits = response.PriceOfThunderFruits;

    var oakFruitsElement = document.getElementById("price-of-oak-fruits");
    var thunderFruitsElement = document.getElementById("price-of-thunder-fruits");
    oakFruitsElement.innerHTML = priceOfOakFruits;
    thunderFruitsElement.innerHTML = priceOfThunderFruits;
}

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

function buyOakFruits() {
    var csrfToken = getCsrfToken();
    var data = {csrfToken : csrfToken};
    params = encodeHTMLForm(data)

    var request = new XMLHttpRequest();
    request.open("POST", "/buy_oak_fruits");
    request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    request.onload = processAfterBuyOakFruits;
    request.send(params);
}

function processAfterBuyOakFruits() {
    var response = JSON.parse(this.response);
    var message = response.Message;
    var money = response.Money;
    var oakFruits = response.OakFruits;
    var thunderFruits = response.ThunderFruits;
    var priceOfOakFruits = response.PriceOfOakFruits;
    var priceOfThunderFruits = response.PriceOfThunderFruits;

    if (message == "success") {
        var moneyElement = document.getElementById("money");
        var oakFruitsElement = document.getElementById("price-of-oak-fruits");
        var thunderFruitsElement = document.getElementById("price-of-thunder-fruits");
        moneyElement.innerHTML = money;
        oakFruitsElement.innerHTML = priceOfOakFruits;
        thunderFruitsElement.innerHTML = priceOfThunderFruits;

        var h6 = document.getElementById("message");

        if (h6 == null) {
            var h6 = document.createElement("h6");
            h6.setAttribute("id", "message");
        }

        h6.innerHTML = "????????????????????????????????????"
        var h5 = document.getElementById("money-h5-element");
        var container = document.getElementById("container");
        container.insertBefore(h6, h5);

    } else if (message == "failed") {
        var h6 = document.getElementById("message");

        if (h6 == null) {
            var h6 = document.createElement("h6");
            h6.setAttribute("id", "message");
        }

        h6.innerHTML = "????????????????????????????????????????????????"
        var h5 = document.getElementById("money-h5-element");
        var container = document.getElementById("container");
        container.insertBefore(h6, h5);

    } else if (message == "not authenticated") {
        window.location.href = "/homepage";
    }
}

function buyThunderFruits() {
    var csrfToken = getCsrfToken();
    var data = {csrfToken : csrfToken};
    params = encodeHTMLForm(data)

    var request = new XMLHttpRequest();
    request.open("POST", "/buy_thunder_fruits");
    request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    request.onload = processAfterBuyThunderFruits;
    request.send(params);
}

function processAfterBuyThunderFruits() {
    var response = JSON.parse(this.response);
    var message = response.Message;
    var money = response.Money;
    var oakFruits = response.OakFruits;
    var thunderFruits = response.ThunderFruits;
    var priceOfOakFruits = response.PriceOfOakFruits;
    var priceOfThunderFruits = response.PriceOfThunderFruits;

    if (message == "success") {
        var moneyElement = document.getElementById("money");
        var oakFruitsElement = document.getElementById("price-of-oak-fruits");
        var thunderFruitsElement = document.getElementById("price-of-thunder-fruits");
        moneyElement.innerHTML = money;
        oakFruitsElement.innerHTML = priceOfOakFruits;
        thunderFruitsElement.innerHTML = priceOfThunderFruits;

        var h6 = document.getElementById("message");

        if (h6 == null) {
            var h6 = document.createElement("h6");
            h6.setAttribute("id", "message");
            var h5 = document.getElementById("money-h5-element");
            var container = document.getElementById("container");
            container.insertBefore(h6, h5);
        }

        h6.innerHTML = "???????????????????????????????????????"

    } else if (message == "failed") {
        var h6 = document.getElementById("message");

        if (h6 == null) {
            var h6 = document.createElement("h6");
            h6.setAttribute("id", "message");
            var h5 = document.getElementById("money-h5-element");
            var container = document.getElementById("container");
            container.insertBefore(h6, h5);
        }

        h6.innerHTML = "???????????????????????????????????????????????????"

    } else if (message == "not authenticated") {
        window.location.href = "/homepage";
    }

}