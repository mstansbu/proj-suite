package tictactoe

import (
	pb "github.com/mstansbu/tic-tac-toe/proto"
)

type TicTacToe struct {
	board [9]byte
}

func NewTicTacToeGame() *TicTacToe {
	return &TicTacToe{board: [9]byte{}}
}

func (ttt *TicTacToe) PlayTurn(payload *pb.Payload) bool {
	ptPayload := payload.GetTttPlayTurnType()
	if ptPayload.FirstPlayer {
		ttt.board[ptPayload.SquarePlayed] = 1
	} else {
		ttt.board[ptPayload.SquarePlayed] = 2
	}
	return ttt.CheckWin()
}

func (ttt *TicTacToe) CheckWin() bool {
	return ttt.checkRows() || ttt.checkColumns() || ttt.checkCross()
}

func (ttt *TicTacToe) checkRows() bool {
	if ttt.board[0] != 0 && ttt.board[0] == ttt.board[1] && ttt.board[0] == ttt.board[2] {
		return true
	}
	if ttt.board[3] != 0 && ttt.board[3] == ttt.board[4] && ttt.board[3] == ttt.board[5] {
		return true
	}
	if ttt.board[6] != 0 && ttt.board[6] == ttt.board[7] && ttt.board[6] == ttt.board[8] {
		return true
	}
	return false
}

func (ttt *TicTacToe) checkColumns() bool {
	if ttt.board[0] != 0 && ttt.board[0] == ttt.board[3] && ttt.board[0] == ttt.board[6] {
		return true
	}
	if ttt.board[1] != 0 && ttt.board[1] == ttt.board[4] && ttt.board[1] == ttt.board[7] {
		return true
	}
	if ttt.board[2] != 0 && ttt.board[2] == ttt.board[5] && ttt.board[2] == ttt.board[8] {
		return true
	}
	return false
}

func (ttt *TicTacToe) checkCross() bool {
	if ttt.board[0] != 0 && ttt.board[0] == ttt.board[4] && ttt.board[0] == ttt.board[8] {
		return true
	}
	if ttt.board[2] != 0 && ttt.board[2] == ttt.board[4] && ttt.board[2] == ttt.board[6] {
		return true
	}
	return false
}
