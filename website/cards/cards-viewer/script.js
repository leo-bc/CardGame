async function loadCards() {
    var cards = await sendGET("/cards/" + gameID);
    console.log(cards);
    var list = document.getElementById("card-list");
    list.innerHTML = '';

    if (cards != null) {
        var prefab = document.getElementById("card-prefab");
        for (var i = 0; i < cards.length; i++) {
            var clone = prefab.cloneNode(true);
            clone.hidden = false;
            clone.dataset.id = cards[i].ID;
            clone.id = `card-${cards[i].ID}`;
            clone.querySelector("#title-text").innerHTML = cards[i].Info.Title;
            clone.querySelector("#card-type-text").innerHTML = cards[i].Info.CardType;
            clone.querySelector("#description-text").innerHTML = cards[i].Info.Description;
            if (cards[i].Info.Ranking == "Legend") {
                clone.style.backgroundColor = '#dbc96e';
            } else if (cards[i].Info.Ranking == "Rare") {
                clone.style.backgroundColor = '#ababab';
            } else {
                clone.style.backgroundColor = '#b8946a';
            }
            list.appendChild(clone);
        }
    }
}

async function createCard() {
    await sendGET("/create-card/");
    loadCards();
}

function backToGame() {
    window.location.href = "/website/games/game-viewer/?id=" + gameID
}


async function removeCard(element) {
    var id = element.dataset.id;
    await sendPOST(`/remove-card/${id}`);
    loadCards();
}

function editCard(element) {
    var id = element.dataset.id;
    window.location.href = `/static/card-studio/card-editor?id=${id}/`;
}

async function setPlayerName() {
    var player = await sendGET("/current-player/");
    document.getElementById("name-text").innerHTML = player.Name;
}


setPlayerName();

const queryString = window.location.search;
const urlParams = new URLSearchParams(queryString);
var gameID = urlParams.get("game-id");
if (gameID != null) {
    loadCards(gameID);
} else {
    document.innerHTML = "euh";
}