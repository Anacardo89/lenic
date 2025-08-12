package wsconnman

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WSConnMan struct {
	clients  map[string]*Client
	mu       sync.Mutex
	Upgrader websocket.Upgrader
}

func NewWSConnMan() *WSConnMan {
	return &WSConnMan{
		clients: make(map[string]*Client),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (cm *WSConnMan) AddClient(username string, conn *websocket.Conn) {

	client := &Client{
		conn:     conn,
		username: username,
		message:  make(chan []byte),
	}
	cm.mu.Lock()
	cm.clients[username] = client
	cm.mu.Unlock()

	go client.listen()
}

func (cm *WSConnMan) RemoveClient(username string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if client, ok := cm.clients[username]; ok {
		close(client.message)
		client.conn.Close()
		delete(cm.clients, username)
	}
}

func (cm *WSConnMan) GetClient(username string) (*Client, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	client, exists := cm.clients[username]
	return client, exists
}

func (cm *WSConnMan) SendMessage(username string, message []byte) error {
	client, exists := cm.GetClient(username)
	if !exists {
		return fmt.Errorf("user %s not connected", username)
	}

	client.message <- message
	return nil
}

func (cm *WSConnMan) IsConnected(username string) bool {
	_, exists := cm.GetClient(username)
	return exists
}
