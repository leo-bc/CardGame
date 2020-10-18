async function loadCards() {
    var cards = await sendGET("/cards/" + gameID);
    console.log(cards);
    var list = document.getElementById("card-list");
    list.innerHTML = '';

    if (cards != null) {
        var attackPrefab = document.getElementById("attack-prefab");
        var cardPrefab = document.getElementById("card-prefab");
        for (var i = 0; i < cards.length; i++) {
            var clone = cardPrefab.cloneNode(true);
            clone.hidden = false;
            clone.dataset.id = cards[i].ID;
            clone.id = `card-${cards[i].ID}`;
            clone.querySelector("#title-text").innerHTML = cards[i].Info.Identity.Title;
            clone.querySelector("#card-type-text").innerHTML = cards[i].Info.Identity.Type;
            clone.querySelector("#description-text").innerHTML = cards[i].Info.Identity.Description;
            if (cards[i].Info.Rank.Ranking == "Legend") {
                clone.style.backgroundColor = '#dbc96e';
            } else if (cards[i].Info.Rank.Ranking == "Rare") {
                clone.style.backgroundColor = '#ababab';
            } else {
                clone.style.backgroundColor = '#b8946a';
            }

            for (var j = 0; j < cards[i].Info.Attacks.length; j++) {
                var attack = cards[i].Info.Attacks[j];
                var attackClone = attackPrefab.cloneNode(true);
                attackClone.hidden = false;
                attackClone.id = `attack-${j}`;
                attackClone.querySelector("#name-text").innerHTML = attack.Name;
                attackClone.querySelector("#cost-text").innerHTML = attack.Cost;
                attackClone.querySelector("#damage-text").innerHTML = attack.Damage;
                clone.querySelector("#attack-list").appendChild(attackClone);
            }

            list.appendChild(clone);
        }
    }
}



function backToGame() {
    window.location.href = "/website/games/game-viewer/?id=" + gameID;
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