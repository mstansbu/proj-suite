"use strict";

let xplayer = true;
let myTurn = xplayer;

let playTurn;

document.body.addEventListener('ttt-start', (event) => {
    const div = document.getElementById("ttt-wrapper");
    xplayer = 'true' === div.dataset.xplayer;
    const gameId = div.dataset.gameid
    myTurn = xplayer;

    const socket = new WebSocket(`ws://localhost:3000/tictactoe/connect/${gameId}`);

    socket.onopen = () => {
        console.log('Game connection established');
    };

    socket.onmessage = (event) => {
        console.log(`Message from server: ${event.data}`);
        const message = JSON.parse(event.data);
        if (message.messageType == "MT_MESSAGEFAIL") {
            console.error("Woop an error from a message");
            //Handle Error
        } else if (message.messageType == "MT_GAMEWIN") {
            console.log("Someone won the game, was it you?");
            myTurn = false;
            //Handle Win
        } else if (message.messageType == "MT_PLAYTURN") {
            const payload = message.payload.tttPlayTurnType
            const id = payload.hasOwnProperty("squarePlayed") ? "ttt-" + payload.squarePlayed : "ttt-0";
            const div = document.getElementById(id);
            if (payload.firstPlayer != xplayer && div.classList.contains('unplayed')) {
                div.classList.remove('unplayed');
                xplayer ? div.classList.add('played-o') : div.classList.add('played-x');
                myTurn = !myTurn;
            }
        }
    };

    socket.onclose = () => {
        console.log('WebSocket connection closed');
    };

    playTurn = (div) => {
        if (myTurn && div.classList.contains('unplayed')) {
            div.classList.remove('unplayed');
            xplayer ? div.classList.add('played-x') : div.classList.add('played-o');
            myTurn = !myTurn;
            const payload = { "messageType": "MT_PLAYTURN", "payload": { "squarePlayed": $(`#${div.id}`).index(), "firstPlayer": xplayer, }, };
            console.log(payload)
            socket.send(JSON.stringify(payload));
        }
    };
});
