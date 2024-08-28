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

func (cm *WSConnManager) AddClient(userID string, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	client := &Client{
		conn:    conn,
		userID:  userID,
		message: make(chan []byte),
	}

	cm.clients[userID] = client

	go client.listen()
}

func (cm *WSConnManager) RemoveClient(userID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if client, ok := cm.clients[userID]; ok {
		close(client.message)
		client.conn.Close()
		delete(cm.clients, userID)
	}
}

func (cm *WSConnManager) GetClient(userID string) (*Client, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	client, exists := cm.clients[userID]
	return client, exists
}

func (cm *WSConnManager) SendMessage(userID string, message []byte) error {
	client, exists := cm.GetClient(userID)
	if !exists {
		return fmt.Errorf("user %s not connected", userID)
	}

	client.message <- message
	return nil
}

func (cm *WSConnManager) IsConnected(userID string) bool {
	_, exists := cm.GetClient(userID)
	return exists
}
