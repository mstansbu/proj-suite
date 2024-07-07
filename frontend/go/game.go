package main

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/valyala/fastjson"
)

const (
	MessageTurnPlayed byte = iota
	MessageGameWin
	MessageFail
)

type GameConnection struct {
	gameServer *GameServer
	Id         uuid.UUID
	board      [9]byte
	players    map[uuid.UUID]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type PlayMessage struct {
	firstPlayer  bool  `json:"firstPlayer"`
	squarePlayed uint8 `json:"squarePlayed"`
}

func NewGameConnection(gs *GameServer) *GameConnection {
	return &GameConnection{
		gameServer: gs,
		Id:         uuid.New(),
		board:      [9]byte{},
		players:    make(map[uuid.UUID]*Client),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
	}
}

func (gc *GameConnection) run() {
	for {
		select {
		case player := <-gc.register:
			gc.players[player.Id] = player
		case player := <-gc.unregister:
			if _, ok := gc.players[player.Id]; ok {
				delete(gc.players, player.Id)
				close(player.send)
			}
		case message := <-gc.broadcast:
			clientId, err := uuid.FromBytes(message[:16])
			if err != nil {
				slog.Error("Something funky happened", "Error", err)
				panic(err)
			}
			payload := message[16:]
			var parser fastjson.Parser
			val, err := parser.ParseBytes(payload)
			if err != nil || !val.Exists("firstPlayer") || !val.Exists("squarePlayed") {
				gc.players[clientId].send <- []byte{MessageFail}
				continue
			}
			Something := val.Get("firstPlayer")
			firstPlayerString := Something.String()
			firstPlayer := true
			if firstPlayerString == "false" {
				firstPlayer = false
			}
			squarePlayed := byte(val.GetInt("squarePlayed"))
			if squarePlayed > 8 {
				gc.players[clientId].send <- []byte{MessageFail}
				continue
			}
			//clientIdByteArray := [16]byte(clientId)
			var firstPlayerByte byte
			if firstPlayer {
				firstPlayerByte = 1
			}
			win := gc.playTurn(firstPlayer, squarePlayed)
			for _, player := range gc.players {
				if clientId != player.Id {
					player.send <- []byte{MessageTurnPlayed, squarePlayed, firstPlayerByte}
				}
				if win {
					//payload := append([]byte{MessageGameWin}, clientIdByteArray[:]...)
					player.send <- []byte{MessageGameWin, firstPlayerByte}
				}
			}
		}
	}
}

func (gc *GameConnection) playTurn(firstPlayer bool, squarePlayed byte) bool {
	if firstPlayer {
		gc.board[squarePlayed] = 1
	} else {
		gc.board[squarePlayed] = 2
	}
	return gc.checkWin()
}

func (gc *GameConnection) checkWin() bool {
	return gc.checkRows() || gc.checkColumns() || gc.checkCross()
}

func (gc *GameConnection) checkRows() bool {
	if gc.board[0] != 0 && gc.board[0] == gc.board[1] && gc.board[0] == gc.board[2] {
		return true
	}
	if gc.board[3] != 0 && gc.board[3] == gc.board[4] && gc.board[3] == gc.board[5] {
		return true
	}
	if gc.board[6] != 0 && gc.board[6] == gc.board[7] && gc.board[6] == gc.board[8] {
		return true
	}
	return false
}

func (gc *GameConnection) checkColumns() bool {
	if gc.board[0] != 0 && gc.board[0] == gc.board[3] && gc.board[0] == gc.board[6] {
		return true
	}
	if gc.board[1] != 0 && gc.board[1] == gc.board[4] && gc.board[1] == gc.board[7] {
		return true
	}
	if gc.board[2] != 0 && gc.board[2] == gc.board[5] && gc.board[2] == gc.board[8] {
		return true
	}
	return false
}

func (gc *GameConnection) checkCross() bool {
	if gc.board[0] != 0 && gc.board[0] == gc.board[5] && gc.board[0] == gc.board[8] {
		return true
	}
	if gc.board[2] != 0 && gc.board[2] == gc.board[5] && gc.board[2] == gc.board[6] {
		return true
	}
	return false
}
