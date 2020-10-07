async function loadPlayers() {
    var players = await sendGET("/players/");
    var list = document.getElementById("player-list");
    list.innerHTML = '';

    if (players != null) {
        var prefab = document.getElementById("player-prefab");
        for (var i = 0; i < players.length; i++) {
            var clone = prefab.cloneNode(true);
            clone.hidden = false;
            clone.dataset.id = players[i].ID;
            clone.id = `card-${players[i].ID}`;
            clone.querySelector("#name-text").innerHTML = players[i].Info.Name;
            var r = Math.floor((players[i].ID * 323.4334) % 64);
            var g = Math.floor((players[i].ID * 123.74) % 64);
            var b = Math.floor((players[i].ID * 523.434) % 64);
            clone.style.backgroundColor = `rgb(${r + 128}, ${g + 128}, ${b + 128})`;
            list.appendChild(clone);
        }
    }
}

async function choosePlayer(element) {
    var id = element.dataset.id;
    await sendPOST(`/select-player/${id}`);
}


loadPlayers();