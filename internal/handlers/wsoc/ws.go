package wsoc

import (
	"encoding/json"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/Anacardo89/tpsi25_blog/pkg/wsocket"
	"github.com/gorilla/websocket"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/ws ", r.RemoteAddr)

	conn, err := wsocket.WSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error.Println("/ws - Failed to upgrade to websocket:", err)
		return
	}

	username := r.URL.Query().Get("user_name")
	if username == "" {
		logger.Error.Println("/ws - No user ID provided", err)
		return
	}

	wsocket.WSConnMan.AddClient(username, conn)
	defer func() {
		wsocket.WSConnMan.RemoveClient(username)
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				logger.Info.Println("/ws - Connection closed normally:", err)
			} else {
				logger.Error.Println("/ws - Could not read message:", err)
			}
			break
		}

		var msg wsocket.Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			logger.Error.Println("/ws - Could not unmarshal message:", err)
			continue
		}

		logger.Info.Printf("/ws - Received message from user %s: %s\n", username, string(message))

		switch msg.Type {
		case "rate_comment":
			handleCommentRate(msg)
		case "rate_post":
			handlePostRate(msg)
		case "comment_on_post":
			handleCommentOnPost(msg)
		case "follow_accept":
			handleFollowAccept(msg)
		case "follow_request":
			handleFollowRequest(msg)
		case "dm":
			handleDM(msg)
		default:
			logger.Warn.Printf("/ws - Unknown message type: %s\n", msg.Type)
		}
	}
}
