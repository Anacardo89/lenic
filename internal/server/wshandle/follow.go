package wshandle

import (
	"encoding/base64"
	"encoding/json"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/repo"
)

func (h *WSHandler) handleFollowRequest(msg Message) {
	// Error Handling
	fail := func(logMsg string, e error) {
		h.log.Error(logMsg, "error", e,
			"message_type", msg.Type,
		)
	}
	//

	// Execution
	bytes, err := base64.URLEncoding.DecodeString(msg.ResourceID)
	if err != nil {
		fail("could not decode user", err)
		return
	}
	username := string(bytes)
	// Early return
	if username == msg.FromUserName {
		return
	}
	// DB operations
	uDB, err := h.db.GetUserByUserName(h.ctx, username)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	fuDB, err := h.db.GetUserByUserName(h.ctx, msg.FromUserName)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	n := &repo.Notification{
		UserID:     uDB.ID,
		FromUserID: fuDB.ID,
		NotifType:  msg.Type,
		NotifText:  msg.Msg,
		ResourceID: msg.ResourceID,
		ParentID:   "",
	}
	if err := h.db.CreateNotification(h.ctx, n); err != nil {
		fail("dberr: could not create notification", err)
		return
	}
	// Response
	u := models.FromDBUserNotif(uDB)
	fu := models.FromDBUserNotif(fuDB)
	notif := models.FromDBNotification(n, *u, *fu)
	notif.ParentID = ""
	data, err := json.Marshal(notif)
	if err != nil {
		fail("failed to marshal JSON", err)
		return
	}
	if h.wsConnMann.IsConnected(u.Username) {
		h.wsConnMann.SendMessage(u.Username, data)
	}
}

func (h *WSHandler) handleFollowAccept(msg Message) {
	// Error Handling
	fail := func(logMsg string, e error) {
		h.log.Error(logMsg, "error", e,
			"message_type", msg.Type,
		)
	}
	//

	// Execution
	bytes, err := base64.URLEncoding.DecodeString(msg.ResourceID)
	if err != nil {
		fail("could not decode user", err)
		return
	}
	username := string(bytes)
	// Early return
	if username == msg.FromUserName {
		return
	}
	// DB operations
	uDB, err := h.db.GetUserByUserName(h.ctx, username)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	fuDB, err := h.db.GetUserByUserName(h.ctx, msg.FromUserName)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	n := &repo.Notification{
		UserID:     uDB.ID,
		FromUserID: fuDB.ID,
		NotifType:  msg.Type,
		NotifText:  msg.Msg,
		ResourceID: msg.ResourceID,
		ParentID:   "",
	}
	if err := h.db.CreateNotification(h.ctx, n); err != nil {
		fail("dberr: could not create notification", err)
		return
	}
	// Response
	u := models.FromDBUserNotif(uDB)
	fu := models.FromDBUserNotif(fuDB)
	notif := models.FromDBNotification(n, *u, *fu)
	notif.ParentID = ""
	data, err := json.Marshal(notif)
	if err != nil {
		fail("failed to marshal JSON", err)
		return
	}
	if h.wsConnMann.IsConnected(u.Username) {
		h.wsConnMann.SendMessage(u.Username, data)
	}
}
