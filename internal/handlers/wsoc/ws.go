package wsoc

import (
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/Anacardo89/tpsi25_blog/pkg/wsocket"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/ws ", r.RemoteAddr)

	conn, err := wsocket.WSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error.Println("/ws - Failed to upgrade to websocket:", err)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		logger.Error.Println("/ws - No user ID provided", err)
		return
	}

	wsocket.WSConnMan.AddClient(userID, conn)
	defer wsocket.WSConnMan.RemoveClient(userID)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Error.Println("/ws - Could not read message", err)
			break
		}
		logger.Info.Printf("/ws - Received message from user %s: %s\n", userID, string(message))
	}
}
