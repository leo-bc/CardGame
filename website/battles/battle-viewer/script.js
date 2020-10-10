async function loadBattle() {
    var side0 = await sendGET(`/battle-side/${gameID}/${battleID}/0`);
    var side1 = await sendGET(`/battle-side/${gameID}/${battleID}/1`);
    document.getElementById("battle-side-list").innerHTML = ""

    showSide(side0);
    showSide(side1);
}

function showSide(side) {
    if (side.IsPlayer) {
        showPlayerSide(side);
    } else {
        showOpponentSide(side);
    }
}

function showPlayerSide(side) {
    var sidePrefab = document.getElementById("player-side-prefab");
    var clone = sidePrefab.cloneNode(true);
    clone.hidden = false;

    var cardSlotPrefab = document.getElementById("card-slot-prefab");
    showBench(side.Info.Bench, clone, cardSlotPrefab);
    showHand(side.Info.Hand, clone, cardSlotPrefab);
    showTakePile(side.Info.TakePile, clone, cardSlotPrefab);

    document.getElementById("battle-side-list").appendChild(clone);
}

function showOpponentSide(side) {
    var sidePrefab = document.getElementById("opponent-side-prefab");
    var clone = sidePrefab.cloneNode(true);
    clone.hidden = false;

    var cardSlotPrefab = document.getElementById("card-slot-prefab");
    showBench(side.Info.Bench, clone, cardSlotPrefab);
    showHand(side.Info.Hand, clone, cardSlotPrefab);
    showTakePile(side.Info.TakePile, clone, cardSlotPrefab);

    document.getElementById("battle-side-list").appendChild(clone);
}

function showBench(bench, parent, cardSlotPrefab) {
    var benchList = parent.querySelector("#bench-list");
    benchList.innerHTML = "";
    for (var i = 0; i < bench.length; i++) {
        var clone = cardSlotPrefab.cloneNode(true);
        clone.hidden = false;
        clone.id = `bench-slot-${i}`;
        clone.querySelector("#card-text").innerHTML = "Card: " + bench[i].CardID;
        benchList.appendChild(clone);
    }
}

function showHand(hand, parent, cardSlotPrefab) {
    var handList = parent.querySelector("#hand-list");
    handList.innerHTML = "";
    if (hand != null) {
        for (var i = 0; i < hand.length; i++) {
            var clone = cardSlotPrefab.cloneNode(true);
            clone.hidden = false;
            clone.id = `hand-slot-${i}`;
            clone.querySelector("#card-text").innerHTML = "Card: " + hand[i].CardID;
            handList.appendChild(clone);
        }
    }
}

function showTakePile(takePile, parent, cardSlotPrefab) {
    var takePileList = parent.querySelector("#take-pile-list")
    takePileList.innerHTML = "";
    if (takePile != null) {
        for (var i = 0; i < takePile.length; i++) {
            var clone = cardSlotPrefab.cloneNode(true);
            clone.hidden = false;
            clone.id = `take-pile-slot-${i}`;
            clone.querySelector("#card-text").innerHTML = "Card: " + takePile[i].CardID;
            takePileList.appendChild(clone);
        }
    }
}


async function drawCard(element) {
    await sendPOST(`/draw-card/${gameID}/${battleID}/`, "");
    loadBattle();
}

const queryString = window.location.search;
const urlParams = new URLSearchParams(queryString);
var gameID = urlParams.get("game-id");
var battleID = urlParams.get("battle-id");
if (gameID != null && battleID != null) {
    loadBattle(gameID, battleID);
}