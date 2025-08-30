package wshandle

import (
	"encoding/base64"
	"encoding/json"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/pkg/logger"
)

func (h *WSHandler) handleFollowRequest(msg Message) {

	bytes, err := base64.URLEncoding.DecodeString(msg.ResourceID)
	if err != nil {
		logger.Error.Printf("Could not decode user %s: %s\n", msg.ResourceID, err)
		return
	}
	userName := string(bytes)

	dbUser, err := h.db.GetUserByUserName(h.ctx, userName)
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

	n := &repo.Notification{
		UserID:     u.ID,
		FromUserID: fromUser.ID,
		NotifType:  msg.Type,
		NotifText:  msg.Msg,
		ResourceID: msg.ResourceID,
		ParentID:   "",
	}

	notifID, err := h.db.CreateNotification(h.ctx, n)
	if err != nil {
		logger.Error.Println("Could not create notification: ", err)
		return
	}

	dbNotif, err := h.db.GetNotification(h.ctx, notifID)
	if err != nil {
		logger.Error.Println("Could not get notification: ", err)
		return
	}
	notif := models.FromDBNotification(dbNotif, *u, *fromU)
	notif.ParentID = ""

	data, err := json.Marshal(notif)
	if err != nil {
		logger.Error.Println("Could not marshal JSON: ", err)
		return
	}

	if h.wsConnMann.IsConnected(dbUser.UserName) {
		h.wsConnMann.SendMessage(u.UserName, data)
	}
}

func (h *WSHandler) handleFollowAccept(msg Message) {

	bytes, err := base64.URLEncoding.DecodeString(msg.ResourceID)
	if err != nil {
		logger.Error.Printf("Could not decode user %s: %s\n", msg.ResourceID, err)
		return
	}
	userName := string(bytes)

	dbUser, err := h.db.GetUserByUserName(h.ctx, userName)
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

	n := &repo.Notification{
		UserID:     u.ID,
		FromUserID: fromUser.ID,
		NotifType:  msg.Type,
		NotifText:  msg.Msg,
		ResourceID: msg.ResourceID,
		ParentID:   "",
	}

	notifID, err := h.db.CreateNotification(h.ctx, n)
	if err != nil {
		logger.Error.Println("Could not create notification: ", err)
		return
	}

	dbNotif, err := h.db.GetNotification(h.ctx, notifID)
	if err != nil {
		logger.Error.Println("Could not get notification: ", err)
		return
	}
	notif := models.FromDBNotification(dbNotif, *u, *fromU)
	notif.ParentID = ""

	data, err := json.Marshal(notif)
	if err != nil {
		logger.Error.Println("Could not marshal JSON: ", err)
		return
	}

	if h.wsConnMann.IsConnected(dbUser.UserName) {
		h.wsConnMann.SendMessage(u.UserName, data)
	}
}
