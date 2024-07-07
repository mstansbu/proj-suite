package main

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	maxMessageSize = 512

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	Id   uuid.UUID
	game *GameConnection
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) read() {
	defer func() {
		c.game.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Connection from client closed unexpectedly", "Error", err)
			}
			break
		}
		message = append(c.Id[:], message...)
		c.game.broadcast <- message
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			var builder strings.Builder
			switch message[0] {
			case MessageFail:
				builder.WriteString(`{"error":"Message Failed To Send"}`)
			case MessageTurnPlayed:
				builder.WriteString(`{"square":` + strconv.Itoa(int(message[1])) + `,"firstPlayer":`)
				if message[2] == 1 {
					builder.WriteString("true")
				} else {
					builder.WriteString("false")
				}
				builder.WriteString(`}`)
			case MessageGameWin:
				builder.WriteString(`{"won":true,"firstPlayer":`)
				if message[1] == 1 {
					builder.WriteString("true")
				} else {
					builder.WriteString("false")
				}
				builder.WriteString(`}`)
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(builder.String()))
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
