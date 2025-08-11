package websockethandle

import (
	"encoding/json"
	"strconv"

	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/pkg/logger"
)

func (h *WSHandler) HandlePostTag(msg Message, taggedUser string) {

	dbUser, err := h.db.GetUserByUserName(h.ctx, taggedUser)
	if err != nil {
		logger.Error.Println("Could not get user: ", err)
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

	n := &db.Notification{
		UserID:     dbUser.ID,
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

func (h *WSHandler) HandleCommentTag(msg Message, taggedUser string) {
	commentID, err := strconv.Atoi(msg.ResourceID)
	if err != nil {
		logger.Error.Printf("Could not convert %s to int: %s\n", msg.ResourceID, err)
		return
	}
	c, err := h.db.GetComment(h.ctx, commentID)
	if err != nil {
		logger.Error.Println("Could not get comment: ", err)
		return
	}
	dbUser, err := h.db.GetUserByUserName(h.ctx, taggedUser)
	if err != nil {
		logger.Error.Println("Could not get user: ", err)
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

	n := &db.Notification{
		UserID:     dbUser.ID,
		FromUserID: fromUser.ID,
		NotifType:  msg.Type,
		NotifText:  msg.Msg,
		ResourceID: msg.ResourceID,
		ParentID:   msg.ParentID,
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
	notif.ParentID = c.PostID

	data, err := json.Marshal(notif)
	if err != nil {
		logger.Error.Println("Could not marshal JSON: ", err)
		return
	}

	if h.wsConnMann.IsConnected(dbUser.UserName) {
		h.wsConnMann.SendMessage(u.UserName, data)
	}
}
