package main

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var ErrGameAlreadyStarted error = errors.New("game has already started")
var ErrGameNotFound error = errors.New("game has not been created yet")
var ErrNoGamesWaiting error = errors.New("no games waiting to be joined")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type GameServer struct {
	gamesInProgress map[uuid.UUID]*GameConnection
	gamesWaiting    map[uuid.UUID]*GameConnection
	register        chan *GameConnection
	unregister      chan *GameConnection
}

func NewGameServer() *GameServer {
	return &GameServer{
		gamesInProgress: make(map[uuid.UUID]*GameConnection),
		gamesWaiting:    make(map[uuid.UUID]*GameConnection),
		register:        make(chan *GameConnection, 256),
		unregister:      make(chan *GameConnection, 256),
	}
}

func (gs *GameServer) run() {
	for {
		select {
		case game := <-gs.register:
			gs.gamesWaiting[game.Id] = game
		case game := <-gs.unregister:
			if _, ok := gs.gamesInProgress[game.Id]; ok {
				delete(gs.gamesInProgress, game.Id)
				close(game.register)
				close(game.unregister)
				close(game.broadcast)
			} else if _, ok := gs.gamesWaiting[game.Id]; ok {
				delete(gs.gamesWaiting, game.Id)
				close(game.register)
				close(game.unregister)
				close(game.broadcast)
			}
		}
	}
}

func (gs *GameServer) startGame(game *GameConnection) error {
	if _, ok := gs.gamesWaiting[game.Id]; !ok {
		if _, ok = gs.gamesInProgress[game.Id]; ok {
			return ErrGameAlreadyStarted
		}
		return ErrGameNotFound
	}
	gs.gamesInProgress[game.Id] = game
	delete(gs.gamesWaiting, game.Id)
	return nil
}

// TODO implement matching logic
func (gs *GameServer) findGame() (*GameConnection, error) {
	if len(gs.gamesWaiting) == 0 {
		return nil, ErrNoGamesWaiting
	}
	for _, game := range gs.gamesWaiting {
		return game, nil
	}
	return nil, ErrNoGamesWaiting
}

func (gs *GameServer) lookUpGame(gameId uuid.UUID) (*GameConnection, error) {
	if game, ok := gs.gamesInProgress[gameId]; ok {
		return game, nil
	} else if game, ok := gs.gamesWaiting[gameId]; ok {
		return game, nil
	}
	return nil, ErrGameNotFound
}
