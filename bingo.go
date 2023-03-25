package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type bingoServer struct {
	serveMux http.ServeMux
	rooms    []*room
	lock     sync.Mutex
}

func (bs *bingoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bs.serveMux.ServeHTTP(w, r)
}

func newBingoServer() *bingoServer {
	bs := new(bingoServer)
	bs.serveMux.Handle("/", http.FileServer(http.Dir("client")))
	bs.serveMux.HandleFunc("/create", bs.createRoomHandler)
	bs.serveMux.HandleFunc("/join", bs.joinRoomHandler)
	return bs
}

type joinRoomRequest struct {
	Username string `json:"username"`
	RoomID   string `json:"roomID"`
}

func (bs *bingoServer) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	newRoom := newRoom()
	bs.lock.Lock()
	bs.rooms = append(bs.rooms, newRoom)
	bs.lock.Unlock()
	res := fmt.Sprintf(`{"roomID":"%s"}`, newRoom.id)
	fmt.Fprintln(w, res)
}

func (bs *bingoServer) joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Print(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "")
	var req joinRoomRequest
	err = wsjson.Read(context.Background(), c, &req)
	if err != nil {
		return
	}
	bs.lock.Lock()
	var roomToAdd *room
	for _, room := range bs.rooms {
		if room.id == req.RoomID {
			roomToAdd = room
		}
	}
	bs.lock.Unlock()
	if roomToAdd == nil {
		c.Close(websocket.StatusProtocolError, "invalid roomID")
		return
	}
	err = roomToAdd.listen(r.Context(), req.Username, c)
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure || websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		log.Print(err)
		return
	}
}
