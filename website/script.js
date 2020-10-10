async function showGame() {
    var gameID = document.getElementById("game-id").value;
    window.location.href = "/website/games/game-viewer/?id=" + gameID;
}

async function loadPlayer() {
    var player = await sendGET("/current-player/");
    if (JSON.stringify(player).length < 3) {
        window.location.href = "/website/players/players-viewer";
    } else {
        document.getElementById("name-text").innerHTML = player.Name;
    }
}

loadPlayer();