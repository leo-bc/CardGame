async function sendPOST(path, body) {
    var json = "{}";
    if (body != null) {
        json = JSON.stringify(body);
    }
    console.log("JSON:", json);

    await fetch(path, {
        method: 'POST', // or 'PUT'
        headers: {
            'Content-Type': 'application/json',
        },
        body: json,
    });
}

async function sendGET(path) {
    try {
        let res = await fetch(path);
        return await res.json();
    } catch (error) {
        return [];
    }
}

function findGetParameter(parameterName) {
    var result = null,
        tmp = [];
    var items = location.search.substr(1).split("&");
    for (var index = 0; index < items.length; index++) {
        tmp = items[index].split("=");
        if (tmp[0] === parameterName) result = decodeURIComponent(tmp[1]);
    }
    return result;
}