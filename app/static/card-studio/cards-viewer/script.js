async function loadCards() {
    var cards = await sendGET("/cards/");
    var list = document.getElementById("card-list");
    list.innerHTML = '';

    if (cards != null) {
        var prefab = document.getElementById("card-prefab");
        for (var i = 0; i < cards.length; i++) {
            var clone = prefab.cloneNode(true);
            clone.hidden = false;
            clone.id = `card-${cards[i].ID}`;
            clone.querySelector("#title-text").innerHTML = cards[i].Info.Title;
            clone.querySelector("#hp-text").innerHTML = cards[i].Info.HP;
            var r = Math.floor((cards[i].ID * 323.4334) % 64);
            var g = Math.floor((cards[i].ID * 123.74) % 64);
            var b = Math.floor((cards[i].ID * 523.434) % 64);
            clone.style.backgroundColor = `rgb(${r + 128}, ${g + 128}, ${b + 128})`;
            list.appendChild(clone);
        }
    }
}

async function createCard() {
    await sendGET("/create-card/");
    loadCards();
}

async function removeCard(element) {
    var id = getID(element);
    await sendPOST(`/remove-card/${id}`);
    loadCards();
}

function getID(element) {
    return parseInt(element.id.split("-")[1]);
}

function editCard(element) {
    var id = getID(element);
    window.location.href = `/static/card-studio/card-editor?id=${id}/`;
}

loadCards();