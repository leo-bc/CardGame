async function loadGame() {
    var game = await sendGET("/game/" + gameID);
    console.log(game);
    if (game != "") {
        document.getElementById("title-text").innerHTML = game.Title;
        var prefab = document.getElementById("player-prefab");
        var list = document.getElementById("player-list");
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
            list.appendChild(clone);
        }
    } else {
        document.getElementById("editor-box").innerHTML = "game DOES NOT EXIST";
    }
}

const queryString = window.location.search;
const urlParams = new URLSearchParams(queryString);
var gameID = urlParams.get("id");
if (gameID != null) {
    loadGame(urlParams.get("id"));
}