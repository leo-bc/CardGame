async function loadBattle() {
    var battle = await sendGET(`/battle/${gameID}/${battleID}/`);
    IDtoCards = battle.Cards;
    if (Object.keys(battle).length !== 0) {
        document.getElementById("start-battle-button").disabled = battle.Info.IsStarted;
        var status = document.getElementById("status-text");
        if (battle.Info.IsStarted && !battle.Info.IsFinished) {
            status.innerHTML = `kuk`;
        } else if (battle.Info.IsFinished) {
            status.innerHTML = `Game ended! won by ${battle.Info.Winner}`;
        } else {
            status.innerHTML = "Game not started!";
        }

        var side0 = battle.Sides[0];
        var side1 = battle.Sides[1];
        document.getElementById("battle-side-list").innerHTML = "";
        showSide(side0, side1);
        showSide(side1, side0);
    }
}

function setCardSlot(element, slot, hidden) {
    if (slot.CardID == -1) {
        element.querySelector("#card-text").innerHTML = "Hidden card";
    } else {

        var card = IDtoCards[slot.CardID]
        element.querySelector("#card-text").innerHTML = card.Identity.Title;
        element.querySelector("#hp-text").innerHTML = `HP: ${card.Rank.HP - slot.DamageTaken} (${(card.Rank.HP - slot.DamageTaken) / card.Rank.HP * 100}%)`;
        var attackPrefab = document.getElementById("attack-prefab");
        var attackList = element.querySelector("#attack-list");
        attackList.innerHTML = "";
        for (var i = 0; i < card.Attacks.length; i++) {
            var attack = card.Attacks[i];
            var clone = attackPrefab.cloneNode(true);
            clone.hidden = false;
            clone.id = `attack-${i}`;
            clone.dataset.id = i;
            clone.querySelector("#name-text").innerHTML = attack.Name;
            clone.querySelector("#damage-text").innerHTML = attack.Damage;
            clone.querySelector("#attack-button").hidden = hidden;
            attackList.appendChild(clone);
        }
    }
}



function showSide(side, opponentSide) {
    if (side.IsPlayer) {
        showPlayerSide(side, opponentSide);
    } else {
        showOpponentSide(side);
    }
}

async function startBattle() {
    await sendPOST(`/battle-start/${gameID}/${battleID}/`, "")
    refresh();
}

function showPlayerSide(side, opponentSide) {
    var sidePrefab = document.getElementById("player-side-prefab");
    var clone = sidePrefab.cloneNode(true);
    clone.hidden = false;

    var cardSlotPrefab = document.getElementById("card-slot-prefab");
    showBench(side.Cards["Bench"], clone, cardSlotPrefab, !side.IsTurn || opponentSide.Cards["Bench"].length == 0);
    showHand(side.Cards["Hand"], clone, cardSlotPrefab, !side.IsTurn || side.Cards["Bench"].length != 0);
    showTakePile(side.Cards["TakePile"], clone);
    showThrowPile(side.Cards["ThrowPile"], clone);
    showPrizes(side.Cards["Prizes"], clone);

    clone.querySelector("#end-turn-button").disabled = !side.IsTurn;

    document.getElementById("battle-side-list").appendChild(clone);
}

function showOpponentSide(side) {
    var sidePrefab = document.getElementById("opponent-side-prefab");
    var clone = sidePrefab.cloneNode(true);
    clone.hidden = false;

    var cardSlotPrefab = document.getElementById("card-slot-prefab");
    showBench(side.Cards["Bench"], clone, cardSlotPrefab, true);
    showHand(side.Cards["Hand"], clone, cardSlotPrefab, true);
    showTakePile(side.Cards["TakePile"], clone);
    showThrowPile(side.Cards["ThrowPile"], clone);
    showPrizes(side.Cards["Prizes"], clone);

    document.getElementById("battle-side-list").appendChild(clone);
}

function showBench(bench, parent, cardSlotPrefab, hidden) {
    var benchList = parent.querySelector("#bench-list");
    benchList.innerHTML = "";
    for (var i = 0; i < bench.length; i++) {
        var clone = cardSlotPrefab.cloneNode(true);
        clone.hidden = false;
        clone.id = `bench-slot-${i}`;
        clone.dataset.id = i;
        setCardSlot(clone, bench[i], hidden);
        benchList.appendChild(clone);
    }
}

function showHand(hand, parent, cardSlotPrefab, hidden) {
    var handList = parent.querySelector("#hand-list");
    handList.innerHTML = "";
    if (hand != null) {
        for (var i = 0; i < hand.length; i++) {
            var clone = cardSlotPrefab.cloneNode(true);
            clone.hidden = false;
            clone.id = `hand-slot-${i}`;
            clone.dataset.id = i;
            clone.querySelector("#play-card-button").hidden = hidden;
            setCardSlot(clone, hand[i], true);
            handList.appendChild(clone);
        }
    }
}

class AttackInfo {
    constructor(source, attack, target) {
        this.Source = source;
        this.Attack = attack;
        this.Target = target;
    }
}

async function attack(element) {
    var sourceID = parseInt(element.parentElement.parentElement.dataset.id);
    var attackID = parseInt(element.dataset.id);
    var info = new AttackInfo(sourceID, attackID, 0);
    await sendPOST(`/battle-attack/${gameID}/${battleID}/`, info)
    refresh();
}

async function playCard(element) {
    var id = element.dataset.id;
    await sendPOST(`/battle-play-card/${gameID}/${battleID}/${id}`, {})
    refresh();
}

function showTakePile(takePile, parent) {
    var takePileList = parent.querySelector("#take-pile-list")
    takePileList.innerHTML = `Take pile size: ${takePile.length}`;
}

function showThrowPile(throwPile, parent) {
    var throwPileList = parent.querySelector("#throw-pile-list")
    throwPileList.innerHTML = `Throw pile size: ${throwPile.length}`;
}

function showPrizes(prizes, parent) {
    var prizesList = parent.querySelector("#prizes-list")
    prizesList.innerHTML = `Prize cards left: ${prizes.length}`;
}

async function endTurn(element) {
    await sendPOST(`/battle-end-turn/${gameID}/${battleID}/`, "");
    refresh();
}

async function refresh() {
    var updateInfo = await sendGET(`/battle-updated/${gameID}/${battleID}/${updateID}`);
    if (updateInfo.IsUpdated == true) {
        updateID = updateInfo.NewID;
        console.log("UPDATING!");
        loadBattle();
    }
}

function backToGame() {
    window.location.href = `/website/games/game-viewer/?id=${gameID}`
}

var IDtoCards = {};
var updateID = -1;

const queryString = window.location.search;
const urlParams = new URLSearchParams(queryString);
var gameID = urlParams.get("game-id");
var battleID = urlParams.get("battle-id");
if (gameID != null && battleID != null) {
    refresh();
    setInterval(refresh, 1000);
}