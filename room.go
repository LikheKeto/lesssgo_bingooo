package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var (
	ErrRoomFull = errors.New("room is full")
)

type player struct {
	name  string
	board [5][5]struct {
		val    int
		marked bool
	}
	conn *websocket.Conn
}

type room struct {
	id          string
	players     []*player
	msg         chan broadcastMessage
	lock        sync.Mutex
	turn        *player
	board       map[int]bool
	gameRunning bool
}

type broadcastMessage struct {
	Type     string      `json:"type"`
	Content  interface{} `json:"content"`
	Sender   *websocket.Conn
	Receiver *websocket.Conn // nil if broadcast
}

func newRoom() *room {
	r := &room{
		id:      uuid.New().String(),
		msg:     make(chan broadcastMessage),
		players: make([]*player, 0),
	}
	go r.broadcast()
	return r
}

func (r *room) addPlayer(name string, c *websocket.Conn) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if len(r.players) >= 2 {
		return ErrRoomFull
	}
	r.players = append(r.players, &player{
		name: name,
		conn: c,
	})
	return nil
}

func (r *room) deletePlayer(name string, c *websocket.Conn) {
	r.lock.Lock()
	defer r.lock.Unlock()
	for i, p := range r.players {
		if p.conn == c {
			r.players = remove(r.players, i)
			log.Printf("User %s left the room %s.\n", name, r.id)
			r.writeAdminMessage("chat", fmt.Sprintf("%s left the room!", name), nil)
			if r.gameRunning {
				r.gameRunning = false
				r.writeAdminMessage("end", "the game has ended!", nil)
			}
			break
		}
	}
}

func (r *room) listen(ctx context.Context, username string, c *websocket.Conn) error {
	err := r.addPlayer(username, c)
	if err != nil {
		return err
	}
	log.Printf("New user %s joined room %s.\n", username, r.id)
	r.writeAdminMessage("chat", fmt.Sprintf("%s joined the room!", username), nil)
	defer r.deletePlayer(username, c)

	for {
		var msg broadcastMessage
		err := wsjson.Read(ctx, c, &msg)
		if err != nil {
			return err
		}
		switch msg.Type {
		case "start":
			err := r.newGame()
			if err != nil {
				r.writeAdminMessage("error", err.Error(), c)
			}
		case "chat":
			msg.Sender = c
			r.msg <- msg
		case "system":
			r.msg <- msg
		case "move":
			move, ok := msg.Content.(float64)
			if !ok {
				r.writeAdminMessage("error", "bad request", c)
			} else {
				err := r.makeMove(c, int(move))
				if err != nil {
					r.writeAdminMessage("error", err.Error(), c)
				}
			}
		}
	}
}

func (r *room) broadcast() {
	for {
		msg := <-r.msg
		if msg.Receiver != nil {
			wsjson.Write(context.Background(), msg.Receiver, msg)
			continue
		}
		for _, p := range r.players {
			if p.conn != msg.Sender {
				err := wsjson.Write(context.Background(), p.conn, msg)
				if err != nil {
					log.Print(err)
				}
			}
			// TODO: see how to handle this error
		}
	}
}

// func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg broadcastMessage) error {
// 	ctx, cancel := context.WithTimeout(ctx, timeout)
// 	defer cancel()
// 	return wsjson.Write(ctx, c, msg)
// }
