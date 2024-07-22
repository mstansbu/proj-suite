package main

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	pb "github.com/mstansbu/tic-tac-toe/proto"
)

type Game interface {
	PlayTurn(*pb.Payload) bool
	CheckWin() bool
}

type GameConnection struct {
	gameServer *GameServer
	Id         uint64
	Game       Game
	players    map[uint32]*Client
	broadcast  chan *pb.Message
	register   chan *Client
	unregister chan *Client
}

func NewGameConnection(gs *GameServer, game Game) *GameConnection {
	return &GameConnection{
		gameServer: gs,
		Id:         rand.Uint64(), //TODO replace with pull/create from DB
		Game:       game,
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
					fmt.Printf("SquarePlayed: %v\n", ptPayload.SquarePlayed)
					if ptPayload.SquarePlayed > 8 {
						//todo send fail message
					}
					win := gc.Game.PlayTurn(message.Payload)
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
