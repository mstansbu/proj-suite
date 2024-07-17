package main

import (
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mstansbu/tic-tac-toe/proto"
	"github.com/mstansbu/tic-tac-toe/templates"
)

var gameServer *GameServer

func main() {
	gameServer = NewGameServer()
	go gameServer.run()

	router := gin.Default()
	router.StaticFile("main.css", "./templates/main.css")
	router.StaticFile("tic-tac-toe.js", "./templates/tic-tac-toe.js")

	router.GET("/", serveHome)
	router.GET("/tictactoe", serveTicTacToe)
	router.GET("/tictactoe/connect/:gameid", clientTIcTacToeConnect)

	router.Run(":3000")

	//router := mux.NewRouter()
}

func serveHome(c *gin.Context) {
	templates.Layout().Render(c.Request.Context(), c.Writer)
}

func serveTicTacToe(c *gin.Context) {
	firstPlayer := true
	var game *GameConnection

	c.Request.ParseForm()
	if c.Request.Form.Has("new") {
		game = NewGameConnection(gameServer)
		go game.run()
		gameServer.register <- game
	} else {
		tempGame, err := gameServer.findGame()
		if err != nil {
			slog.Error("Error trying to connect to game in progress", "Error", err)
			c.Status(http.StatusNotFound)
			return
		}
		// this is stupid but done because variable scope results in overriding game if not using intermediate variable
		game = tempGame
		firstPlayer = false
		err = gameServer.startGame(game)
		if err != nil {
			slog.Error("Could not start game", "Error", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.Status(http.StatusOK)
	templates.TicTacToe(game.Id, firstPlayer).Render(c.Request.Context(), c.Writer)
}

func clientTIcTacToeConnect(c *gin.Context) {
	var gameId uint64
	gameIdString, ok := c.Params.Get("gameid")
	if !ok {
		slog.Error("Error getting game id param, param not found")
		c.Status(http.StatusBadRequest)
		return
	}
	gameId, err := strconv.ParseUint(gameIdString, 10, 64)
	if err != nil {
		slog.Error("Error parsing game id from url", "Error", err)
		c.Status(http.StatusBadRequest)
		return
	}
	game, err := gameServer.lookUpGame(gameId)
	if err != nil {
		slog.Error("Could not find game in server", "Error", err)
		c.Status(http.StatusBadRequest)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("Failed to upgrade connection", "Error", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	player := &Client{Id: rand.Uint32(), conn: conn, game: game, send: make(chan *proto.Message, 256)}

	game.register <- player

	go player.read()
	go player.write()

	c.Status(http.StatusOK)
}
