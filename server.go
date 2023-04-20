package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader     = websocket.Upgrader{}
	wsConnection *websocket.Conn
	clipContent  string
)

func main() {
	r := gin.Default()
	r.GET("/read", readHandler)
	r.POST("/write", writeHandler)

	r.GET("/ws", websocketHandler)

	err := r.Run(":8080")
	if err != nil {
		log.Println(err)
	}
}

func readHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"content": clipContent})
}

func writeHandler(c *gin.Context) {
	var payload struct {
		Content string `json:"content"`
	}
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid payload",
		})
		log.Println(err)
		return
	}

	clipContent = payload.Content
	sendWsMessage()

	c.JSON(http.StatusOK, gin.H{"message": "content updated"})
}

func websocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
	}
	wsConnection = conn
	go recveveWsMessage()
}

func sendWsMessage() {
	if wsConnection == nil {
		return
	}
	err := wsConnection.WriteMessage(websocket.TextMessage, []byte(clipContent))
	if err != nil {
		log.Println(err)
	}
}

func recveveWsMessage() {
	for {
		if wsConnection == nil {
			return
		}
		_, msg, err := wsConnection.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		clipContent = string(msg)
	}
}
