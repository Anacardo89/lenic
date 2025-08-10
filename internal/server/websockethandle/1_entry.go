package websockethandle

import (
	"encoding/json"
	"net/http"

	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/gorilla/websocket"
)

type Message struct {
	FromUserName string `json:"from_username"`
	Type         string `json:"type"`
	Msg          string `json:"msg"`
	ResourceID   string `json:"resource_id"`
	ParentID     string `json:"parent_id"`
}

func (h *WSHandler) HandleWSMsg(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/ws ", r.RemoteAddr)

	conn, err := h.wsConnMann.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error.Println("/ws - Failed to upgrade to websocket:", err)
		return
	}

	username := r.URL.Query().Get("user_name")
	if username == "" {
		logger.Error.Println("/ws - No user ID provided", err)
		return
	}

	h.wsConnMann.AddClient(username, conn)
	defer func() {
		h.wsConnMann.RemoveClient(username)
		conn.Close()
	}()

	for {

		select {
		case <-h.ctx.Done():
			return
		default:
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				logger.Info.Println("/ws - Connection closed normally:", err)
			} else {
				logger.Error.Println("/ws - Could not read message:", err)
			}
			break
		}

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			logger.Error.Println("/ws - Could not unmarshal message:", err)
			continue
		}

		logger.Info.Printf("/ws - Received message from user %s: %s\n", username, string(message))

		switch msg.Type {
		case "rate_comment":
			h.handleCommentRate(msg)
		case "rate_post":
			h.handlePostRate(msg)
		case "comment_on_post":
			h.handleCommentOnPost(msg)
		case "follow_accept":
			h.handleFollowAccept(msg)
		case "follow_request":
			h.handleFollowRequest(msg)
		case "dm":
			h.handleDM(msg)
		default:
			logger.Warn.Printf("/ws - Unknown message type: %s\n", msg.Type)
		}
	}
}
