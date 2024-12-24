package server

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RoomMap struct {
	mu    sync.Mutex
	Rooms map[string][]User
}

type User struct {
	host bool
	Conn *websocket.Conn
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

func (r *RoomMap) CreateRoom() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	rand.Seed(time.Now().UnixNano())

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 8)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	roomId := string(b)
	r.Rooms[roomId] = []User{}
	return roomId

}

func (r *RoomMap) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := User{host, conn}

	log.Println("Inserting into Room with RoomID: ", roomID)
	r.Rooms[roomID] = append(r.Rooms[roomID], p)

	fmt.Println(r.Rooms)

}

func (r *RoomMap) DeleteRoom(roomID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.Rooms, roomID)
}
