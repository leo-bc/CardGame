async function loadCard() {
    var card = await sendGET("/card/" + cardID);
    console.log(card);
    if (card != "") {
        document.getElementById("title-input").value = card.Title;
        document.getElementById("hp-input").value = card.HP;
    } else {
        document.getElementById("editor-box").innerHTML = "CARD DOES NOT EXIST";
    }
}

async function saveCard() {
    var title = document.getElementById("title-input").value;
    var hp = parseInt(document.getElementById("hp-input").value, 10);
    await sendPOST("/card/" + cardID, new CardInfo(title, hp));
}

async function createCard() {
    var index = await sendGET("/create-card/");
    window.location.href = `/static/card-studio/card-editor?id=${index}/`;
}

function toOverview() {
    window.location.href = `/static/card-studio/cards-viewer/`;
}

const queryString = window.location.search;
const urlParams = new URLSearchParams(queryString);
var cardID = urlParams.get("id");
if (cardID != null) {
    loadCard(urlParams.get("id"));
}