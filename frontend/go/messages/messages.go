package messages

import (
	"errors"

	"github.com/google/uuid"
)

var ErrPayloadMarshall error = errors.New("unable to convert payload from bytes to struct")
var ErrInvalidPayload error = errors.New("payload does not match MessageType")

type MessageType byte

const (
	MessageFail MessageType = iota
	MessageGameWin
	MessageTurnPlayed
	MessagePlayerJoin
	MessageGameStart
)

type Message struct {
	MessageId   uuid.UUID
	MessageType MessageType
	SenderId    uuid.UUID
	Payload     Payload
}

func (message *Message) Validate() bool {
	return message.Payload.checkValidPayload(message.MessageType)
}

type Payload interface {
	checkValidPayload(MessageType) bool
}

type TTTPlayPayload struct {
	FirstPlayer  bool  `json:"firstPlayer"`
	SquarePlayed uint8 `json:"squarePlayed"`
}

func (payload *TTTPlayPayload) checkValidPayload(mType MessageType) bool {
	return mType == MessageTurnPlayed
}

func newTTTPlayPayload(unformattedPayload []byte) *TTTPlayPayload {
	payload := &TTTPlayPayload{FirstPlayer: true, SquarePlayed: unformattedPayload[1]}
	if unformattedPayload[0] == 0 {
		payload.FirstPlayer = false
	}
	return payload
}

func ConstructPayload(mType MessageType, unformattedPayload []byte) (Payload, error) {
	payload := newTTTPlayPayload(unformattedPayload)
	return payload, nil
}
