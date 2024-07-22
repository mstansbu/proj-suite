package main

import (
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	pb "github.com/mstansbu/tic-tac-toe/proto"
	"google.golang.org/protobuf/encoding/protojson"
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

type Client struct {
	Id   uint32
	game *GameConnection
	conn *websocket.Conn
	send chan *pb.Message
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
		_, byteMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("Connection from client closed unexpectedly", "Error", err)
			}
			break
		}
		if len(byteMessage) != 0 {
			jsonMessage, err := ClientParser.ParseBytes(byteMessage)
			if err != nil {
				//TODO
				slog.Error("Parser Error", "Error", err, "byteMessage", byteMessage)
				panic("yurp")
			}
			messageId := uuid.New()
			message := &pb.Message{Id: messageId[:], SenderId: c.Id, ServerId: c.game.Id}
			messageType := jsonMessage.GetStringBytes("messageType")
			payload := jsonMessage.Get("payload")
			switch string(messageType) {
			case "MT_PLAYTURN":
				message.MessageType = pb.Message_MT_PLAYTURN
				message.Payload = &pb.Payload{
					Type: &pb.Payload_TttPlayTurnType{
						TttPlayTurnType: &pb.PayloadPlayTurn{
							FirstPlayer:  payload.GetBool("firstPlayer"),
							SquarePlayed: uint32(payload.GetUint("squarePlayed")),
						}}}
			}
			c.game.broadcast <- message
		}
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
			var toClient []byte
			switch message.MessageType {
			case pb.Message_MT_MESSAGEFAIL:
				//TODO handle message type FAIL
				builder.WriteString(`{"error":"Message Failed To Send"}`)
				toClient = []byte(builder.String())
			case pb.Message_MT_PLAYTURN:
				switch message.Payload.Type.(type) {
				case *pb.Payload_TttPlayTurnType:
					out, err := protojson.Marshal(message)
					if err != nil {
						//TODO handle marshall error
						panic("yurp")
					}
					toClient = out
				case nil:
					//TODO deal with message but no payload
				default:
					//TODO deal with unrecognized payload type
				}
			case pb.Message_MT_GAMEWIN:
				switch message.Payload.Type.(type) {
				case *pb.Payload_TttGameWinType:
					out, err := protojson.Marshal(message)
					if err != nil {
						//TODO handle marshall error
						panic("yurp")
					}
					toClient = out
				case nil:
					//TODO deal with message but no payload
				default:
					//TODO deal with unrecognized payload type
				}
			default:
				//TODO with unrecognize message type
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(toClient)
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
