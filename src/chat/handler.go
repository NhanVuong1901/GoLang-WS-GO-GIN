package chat

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var UPGRADER = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServerWS(c *gin.Context) {
	roomID := c.Query("room")
	userID := c.Query("user")

	conn, err := UPGRADER.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Websocket upgrade failed", err)
		return
	}

	client := &Client{
		Conn:   conn,
		UserID: userID,
		RoomID: roomID,
		Send:   make(chan []byte),
	}

	WS.Register <- client

	go client.ReadPump()
	go client.WritePump()
}

func (c *Client) ReadPump() {
	defer func() {
		WS.UnRegister <- c
		c.Conn.Close()
	}()
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		WS.Broadcast <- &MessagePayload{
			RoomID:  c.RoomID,
			Message: msg,
		}
	}
}
func (c *Client) WritePump() {
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
