package main

import (
	"math/rand"

	"github.com/google/uuid"
	pb "github.com/mstansbu/tic-tac-toe/proto"
)

type GameConnection struct {
	gameServer *GameServer
	Id         uint64
	board      [9]byte
	players    map[uint32]*Client
	broadcast  chan *pb.Message
	register   chan *Client
	unregister chan *Client
}

func NewGameConnection(gs *GameServer) *GameConnection {
	return &GameConnection{
		gameServer: gs,
		Id:         rand.Uint64(), //TODO replace with pull/create from DB
		board:      [9]byte{},
		players:    make(map[uint32]*Client),
		broadcast:  make(chan *pb.Message, 256),
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
			switch message.MessageType {
			case pb.Message_MT_PLAYTURN:
				switch message.Payload.Type.(type) {
				case *pb.Payload_TttPlayTurnType:
					ptPayload := message.Payload.GetTttPlayTurnType()
					if ptPayload.SquarePlayed > 8 {
						//todo send fail message
					}
					win := gc.playTurn(message.Payload.GetTttPlayTurnType().FirstPlayer, message.Payload.GetTttPlayTurnType().GetSquarePlayed())
					var winMessage pb.Message
					for _, player := range gc.players {
						if message.SenderId != player.Id {
							player.send <- message
						}
						if win {
							id := uuid.New() //TODO Figure out message ids and whether they should be one per event or message sent
							winMessage = pb.Message{
								MessageType: pb.Message_MT_GAMEWIN,
								Id:          id[:],
								SenderId:    message.SenderId,
								ServerId:    message.ServerId,
								Payload:     &pb.Payload{Type: &pb.Payload_TttGameWinType{TttGameWinType: &pb.PayloadGameWin{FirstPlayer: ptPayload.FirstPlayer}}},
							}
							player.send <- &winMessage
						}
					}
				case nil:
				default:
				}
			}

		}
	}
}

func (gc *GameConnection) playTurn(firstPlayer bool, squarePlayed uint32) bool {
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
	if gc.board[0] != 0 && gc.board[0] == gc.board[4] && gc.board[0] == gc.board[8] {
		return true
	}
	if gc.board[2] != 0 && gc.board[2] == gc.board[4] && gc.board[2] == gc.board[6] {
		return true
	}
	return false
}
