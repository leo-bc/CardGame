async function loadGame() {
    var game = await sendGET("/game/" + gameID);
    console.log(game);
    if (game != "") {
        document.getElementById("title-input").value = game.Title;
        for (var i = 0; i < game.Players.length; i++) {
            var player = await sendGET("/player/" + game.Players[i].PlayerID);
            document.getElementById("bitch-input").innerHTML += player.Name + ":</br>";
            for (var j = 0; j < game.Players[i].CardIDs.length; j++) {
                var card = await sendGET("/card/" + game.Players[i].CardIDs[j]);
                console.log(card);
                document.getElementById("bitch-input").innerHTML += "Card: " + card.Title + "</br>";
            }
            document.getElementById("bitch-input").innerHTML += "</br>";
        }
        // document.getElementById("hp-input").value = card.HP;
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