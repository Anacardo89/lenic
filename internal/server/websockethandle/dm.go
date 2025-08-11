package websockethandle

import (
	"encoding/json"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/pkg/logger"
)

func (h *WSHandler) handleDM(msg Message) {
	logger.Info.Println("/ws handling DM")
	logger.Debug.Println(msg)

	dbUser, err := h.db.GetUserByUserName(h.ctx, msg.ResourceID)
	if err != nil {
		logger.Error.Println("Could not get post: ", err)
		return
	}

	fromUser, err := h.db.GetUserByUserName(h.ctx, msg.FromUserName)
	if err != nil {
		logger.Error.Println("Could not get from user: ", err)
		return
	}

	if dbUser.ID == fromUser.ID {
		return
	}

	u := models.FromDBUserNotif(dbUser)
	fromU := models.FromDBUserNotif(fromUser)

	dbConvo, err := h.db.GetConversationByUsers(h.ctx, u.ID, fromU.ID)
	if err != nil {
		logger.Error.Println("Could not get conversation: ", err)
		return
	}

	err = h.db.UpdateConversation(h.ctx, dbConvo.ID)
	if err != nil {
		logger.Error.Println("Could not update conversation: ", err)
		return
	}

	n := &models.Notification{
		User:       *u,
		FromUser:   *fromU,
		NotifType:  models.NotifType(msg.Type),
		NotifText:  msg.Msg,
		ResourceID: dbConvo.ID.String(),
		ParentID:   "",
		IsRead:     false,
	}

	data, err := json.Marshal(n)
	if err != nil {
		logger.Error.Println("Could not marshal JSON: ", err)
		return
	}

	if h.wsConnMann.IsConnected(dbUser.UserName) {
		h.wsConnMann.SendMessage(u.UserName, data)
	}
}
