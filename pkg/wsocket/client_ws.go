package wsocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	userID  string
	message chan []byte
}

func (client *Client) listen() {
	defer func() {
		client.conn.Close()
	}()

	for msg := range client.message {
		err := client.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Failed to send message:", err)
			return
		}
	}
}
