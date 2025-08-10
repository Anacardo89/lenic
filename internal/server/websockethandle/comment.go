package websockethandle

import (
	"encoding/json"
	"strconv"

	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/pkg/logger"
)

func (h *WSHandler) handleCommentOnPost(msg Message) {
	comment_id, err := strconv.Atoi(msg.ResourceID)
	if err != nil {
		logger.Error.Printf("Could not convert %s to int: %s\n", msg.ResourceID, err)
		return
	}
	c, err := h.db.GetComment(h.ctx, comment_id)
	if err != nil {
		logger.Error.Println("Could not get comment: ", err)
		return
	}
	dbpost, err := h.db.GetPost(h.ctx, c.PostID)
	if err != nil {
		logger.Error.Println("Could not get post: ", err)
		return
	}
	dbuser, err := h.db.GetUserByID(h.ctx, dbpost.AuthorID)
	if err != nil {
		logger.Error.Println("Could not get user: ", err)
		return
	}

	fromuser, err := h.db.GetUserByUserName(h.ctx, msg.FromUserName)
	if err != nil {
		logger.Error.Println("Could not get from user: ", err)
		return
	}

	if dbuser.ID == fromuser.ID {
		return
	}

	u := models.FromDBUserNotif(dbuser)
	from_u := models.FromDBUserNotif(fromuser)

	n := &db.Notification{
		UserID:     dbpost.AuthorID,
		FromUserID: fromuser.ID,
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

	dbnotif, err := h.db.GetNotification(h.ctx, notifID)
	if err != nil {
		logger.Error.Println("Could not get notification: ", err)
		return
	}
	notif := models.FromDBNotification(dbnotif, *u, *from_u)
	notif.ParentId = c.PostID

	data, err := json.Marshal(notif)
	if err != nil {
		logger.Error.Println("Could not marshal JSON: ", err)
		return
	}

	h.wsConnMann.SendMessage(u.UserName, data)
}
