package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var AllRooms RoomMap

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *RoomMap) Init() {
	r.Rooms = make(map[string][]User)
}

// Rendering home page.
func HomePage(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func CreateRoomHandler(c *gin.Context) {

	roomId := AllRooms.CreateRoom()
	c.HTML(200, "newchat.html", gin.H{
		"roomID": roomId,
	})
}

func JoinRoomHandler(c *gin.Context) {
	roomID := c.PostForm("roomID")

	c.HTML(200, "joinchat.html", gin.H{
		"roomID": roomID,
	})

}
func JoinRoomRequest(c *gin.Context) {

	roomId := c.Query("roomID")
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("Webscoket upgrade error : ", err)

	}

	AllRooms.InsertIntoRoom(roomId, false, ws)
	go Broadcaster()

	for {
		var msg broadcastMsg
		err = ws.ReadJSON(&msg.Message)
		if err != nil {
			fmt.Println("Error getting reading the message from the websocket :", err)
			return
		}
		msg.Client = ws
		msg.RoomID = roomId

		broadcast <- msg
		time.Sleep(1 * time.Second)
	}
}

func Broadcaster() {

	for {
		msg := <-broadcast

		for _, client := range AllRooms.Rooms[msg.RoomID] {
			if client.Conn != msg.Client {

				// Panic Recover
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered from panic:", r)
						client.Conn.Close()
					}
				}()

				err := client.Conn.WriteJSON(msg.Message)
				if err != nil {
					fmt.Println("Error gettign writing the message : ", err)
					client.Conn.Close()
				}
			}
		}

	}
}

func CloseRoomRequestHandler(c *gin.Context) {
	roomID := c.Query("roomID")

	AllRooms.DeleteRoom(roomID)

	c.Redirect(http.StatusSeeOther, "/")
}
