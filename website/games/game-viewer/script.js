async function showGame() {
    var game = await sendGET("/game/" + gameID);
    var list = document.getElementById("player-list");
    list.innerHTML = "";
    if (game != "") {
        document.getElementById("title-text").innerHTML = game.Title;
        var prefab = document.getElementById("player-prefab");

        document.getElementById("join-game-button").disabled = false;
        document.getElementById("leave-game-button").disabled = true;


        for (var i = 0; i < game.Players.length; i++) {
            var playerID = game.Players[i].PlayerID;
            var player = await sendGET("/player/" + playerID);

            var clone = prefab.cloneNode(true);
            clone.hidden = false;
            clone.dataset.id = playerID;
            clone.id = `player-${playerID}`;
            clone.querySelector("#name-text").innerHTML = player.Name;
            var r = Math.floor((playerID * 323.4334) % 64);
            var g = Math.floor((playerID * 123.74) % 64);
            var b = Math.floor((playerID * 523.434) % 64);
            clone.style.backgroundColor = `rgb(${r + 128}, ${g + 128}, ${b + 128})`;

            var readyText = clone.querySelector("#ready-text");
            if (!game.IsStarted) {
                if (game.Players[i].IsReady) {
                    readyText.innerHTML = "Ready!!!!";
                } else {
                    readyText.innerHTML = "Not ready 🥱";
                }

                if (player.Name == currentPlayer.Name) {
                    document.getElementById("join-game-button").disabled = true;
                    document.getElementById("leave-game-button").disabled = false;
                    var button = clone.querySelector("#ready-button");
                    button.hidden = false;
                    if (game.Players[i].IsReady)
                        button.innerHTML = "Unready"
                    else
                        button.innerHTML = "Ready"
                }
            } else {
                readyText.innerHTML = "Game has started";
                document.getElementById("join-game-button").hidden = true;
                document.getElementById("leave-game-button").hidden = true;
            }

            list.appendChild(clone);
        }
    } else {
        document.getElementById("editor-box").innerHTML = "game DOES NOT EXIST";
    }
}

async function setReady() {
    await sendPOST("/set-ready/" + gameID);
    await refresh();
}

async function loadCurrentPlayer() {
    var player = await sendGET("/current-player/");
    if (JSON.stringify(player).length < 3) {
        return null;
    }
    return player;
}

async function startPage() {
    currentPlayer = await loadCurrentPlayer();

    if (currentPlayer != null) {
        const queryString = window.location.search;
        const urlParams = new URLSearchParams(queryString);
        gameID = urlParams.get("id");
        if (gameID != null) {
            showGame(gameID);
        }

        setInterval(refresh, 1000);
    } else {
        window.location.href = "/website/";
    }
}

var currentPlayer;
var gameID;

startPage();

async function refresh() {
    var updated = await sendGET("/game-updated/" + gameID);
    if (updated == true) {
        console.log("UPDATING!");
        showGame();
    }
}

async function joinGame() {
    await sendPOST("/join-game/" + gameID);
    await refresh();
}

async function leaveGame() {
    await sendPOST("/leave-game/" + gameID);
    await refresh();
}