package wshandle

import (
	"encoding/json"
	"net/http"

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
	conn, err := h.wsConnMann.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error("failed to upgrade conn to websocket:", err)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		h.log.Error("no user provided", err)
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
				h.log.Info("connection closed normally:", err)
			} else {
				h.log.Error("could not read message:", err)
			}
			break
		}

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			h.log.Error("could not unmarshal message:", err)
			continue
		}

		h.log.Info("received message from user %s: %s\n", username, string(message))

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
			h.log.Warn("unknown message type:", msg.Type)
		}
	}
}
