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
        const payload = JSON.parse(event.data);
        if (payload.hasOwnProperty("error")) {
            console.error("Woop an error from a message");
            //Handle Error
        } else if (payload.hasOwnProperty("won")) {
            console.log("Someone won the game, was it you?");
            myTurn = false;
            //Handle Win
        } else if (payload.hasOwnProperty("square")) {
            if (payload.firstPlayer != xplayer && div.classList.contains('unplayed')) {
                const id = "ttt-" + payload.square;
                const div = document.getElementById(id);
                div.classList.remove('unplayed');
                xplayer ? div.classList.add('played-o') : div.classList.add('played-x');
                console.log(div.classList)
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
            const payload = { "squarePlayed": $('li').index(div), "firstPlayer": xplayer };
            console.log(payload)
            socket.send(payload);
        }
    };
});
