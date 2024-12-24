package main

import (
	"video_chat/server"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	server.AllRooms.Init()
	router.LoadHTMLGlob("client/*.html")
	router.Static("/static", "client/static")

	router.GET("/", server.HomePage)
	router.GET("/create", server.CreateRoomHandler)
	router.GET("/join", server.JoinRoomRequest)
	router.POST("/joinroom", server.JoinRoomHandler)
	router.GET("/close", server.CloseRoomRequestHandler)

	router.Run(":3500")

}
