package main

import (
	"errors"
	"fmt"

	"nhooyr.io/websocket"
)

func (r *room) writeAdminMessage(msgType string, message string, receiver *websocket.Conn) {
	r.msg <- broadcastMessage{
		Type: msgType,
		Content: map[string]string{
			"message": message,
			"author":  "Admin",
		},
		Receiver: receiver,
	}
}

func (r *room) newGame() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.gameRunning {
		return errors.New("game is already running")
	}
	if len(r.players) != 2 {
		return errors.New("insufficient players to start the game")
	}
	r.gameRunning = true
	r.board = make(map[int]bool)
	for _, player := range r.players {
		uniqBoard := generateRandomBoard()
		for i, row := range player.board {
			for j := range row {
				player.board[i][j] = struct {
					val    int
					marked bool
				}{
					val:    uniqBoard[i][j],
					marked: false,
				}
			}
		}
		r.msg <- broadcastMessage{
			Type:     "board",
			Content:  uniqBoard,
			Receiver: player.conn,
		}
	}
	r.turn = r.players[0]
	r.writeAdminMessage("start", "game has started!", nil)
	r.writeAdminMessage("turn", "It is your turn!", r.turn.conn)
	r.writeAdminMessage("turn", fmt.Sprintf("It is %s's turn!", r.turn.name), r.players[1].conn)
	return nil
}

func (r *room) makeMove(c *websocket.Conn, move int) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if !r.gameRunning {
		return errors.New("the game is not running")
	}
	if r.board[move] {
		return errors.New("the number has already been marked")
	}
	if r.turn.conn != c {
		return errors.New("it is't your turn to make move")
	}
	r.board[move] = true
	r.msg <- broadcastMessage{
		Type:    "move",
		Content: move,
	}

	for _, player := range r.players {
		if player.conn != c {
			r.turn = player
			r.writeAdminMessage("turn", "It is your turn!", r.turn.conn)
			for _, p := range r.players {
				if p != player {
					r.writeAdminMessage("turn", fmt.Sprintf("It is %s's turn!", r.turn.name), p.conn)
				}
			}
		}
	playerloop:
		for i, row := range player.board {
			for j := range row {
				if player.board[i][j].val == move {
					player.board[i][j].marked = true
					break playerloop
				}
			}
		}
		if player.isBingo() {
			r.declareWinner(player)
			r.gameRunning = false
			r.writeAdminMessage("end", "the game has ended!", nil)
		}
	}
	return nil
}

func (r *room) declareWinner(p *player) {
	r.writeAdminMessage("win", "YOU WON!!!", p.conn)
	var loser *player
	for _, player := range r.players {
		if player != p {
			loser = player
		}
	}
	r.writeAdminMessage("loss", "YOU LOST!", loser.conn)
}

func (p *player) isBingo() bool {
	counter := 0
	for i := 0; i < 5; i++ {
		subCounter := 0
		for j := 0; j < 5; j++ {
			if !p.board[i][j].marked {
				break
			}
			subCounter++
		}
		if subCounter == 5 {
			counter++
		}
		subCounter = 0
		for j := 0; j < 5; j++ {
			if !p.board[j][i].marked {
				break
			}
			subCounter++
		}
		if subCounter == 5 {
			counter++
		}
	}
	subCounter := 0
	for i := 0; i < 5; i++ {
		if !p.board[i][i].marked {
			break
		}
		subCounter++
	}
	if subCounter == 5 {
		counter++
	}
	subCounter = 0
	for i := 0; i < 5; i++ {
		if !p.board[4-i][i].marked {
			break
		}
		subCounter++
	}
	if subCounter == 5 {
		counter++
	}
	if counter >= 5 {
		return true
	}
	return false
}
