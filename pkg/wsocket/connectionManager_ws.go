package wsocket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	WSConnMan  *WSConnManager
	WSUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow connections from any origin for simplicity
			return true
		},
	}
)

type WSConnManager struct {
	clients map[string]*Client
	mu      sync.Mutex
}

func NewWSConnManager() *WSConnManager {
	return &WSConnManager{
		clients: make(map[string]*Client),
	}
}

func (cm *WSConnManager) AddClient(username string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	client := &Client{
		conn:     conn,
		username: username,
		message:  make(chan []byte),
	}

	cm.clients[username] = client

	go client.listen()
}

func (cm *WSConnManager) RemoveClient(username string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if client, ok := cm.clients[username]; ok {
		close(client.message)
		client.conn.Close()
		delete(cm.clients, username)
	}
}

func (cm *WSConnManager) GetClient(username string) (*Client, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	client, exists := cm.clients[username]
	return client, exists
}

func (cm *WSConnManager) SendMessage(username string, message []byte) error {
	client, exists := cm.GetClient(username)
	if !exists {
		return fmt.Errorf("user %s not connected", username)
	}

	client.message <- message
	return nil
}

func (cm *WSConnManager) IsConnected(username string) bool {
	_, exists := cm.GetClient(username)
	return exists
}
