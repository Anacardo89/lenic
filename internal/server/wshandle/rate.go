package wshandle

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/repo"
)

func (h *WSHandler) handlePostRate(msg Message) {
	// Error Handling
	fail := func(logMsg string, e error) {
		h.log.Error(logMsg, "error", e,
			"message_type", msg.Type,
		)
	}
	//

	// Execution
	pID, err := uuid.Parse(msg.ResourceID)
	if err != nil {
		fail("parsing post uuid", err)
		return
	}
	// DB operations
	pDB, err := h.db.GetPost(h.ctx, pID)
	if err != nil {
		fail("dberr: could not get post", err)
		return
	}
	uDB, err := h.db.GetUserByID(h.ctx, pDB.AuthorID)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	fuDB, err := h.db.GetUserByUserName(h.ctx, msg.FromUserName)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	if uDB.ID == fuDB.ID {
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

func (h *WSHandler) handleCommentRate(msg Message) {
	// Error Handling
	fail := func(logMsg string, e error) {
		h.log.Error(logMsg, "error", e,
			"message_type", msg.Type,
		)
	}
	//

	// Execution
	cID, err := uuid.Parse(msg.ResourceID)
	if err != nil {
		fail("parsing comment uuid", err)
		return
	}
	cDB, err := h.db.GetComment(h.ctx, cID)
	if err != nil {
		fail("dberr: could not get comment", err)
		return
	}
	uDB, err := h.db.GetUserByID(h.ctx, cDB.AuthorID)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	fuDB, err := h.db.GetUserByUserName(h.ctx, msg.FromUserName)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	if uDB.ID == fuDB.ID {
		return
	}
	n := &repo.Notification{
		UserID:     uDB.ID,
		FromUserID: fuDB.ID,
		NotifType:  msg.Type,
		NotifText:  msg.Msg,
		ResourceID: msg.ResourceID,
		ParentID:   msg.ParentID,
	}
	if err := h.db.CreateNotification(h.ctx, n); err != nil {
		fail("dberr: could not create notification", err)
		return
	}
	// Response
	u := models.FromDBUserNotif(uDB)
	fu := models.FromDBUserNotif(fuDB)
	notif := models.FromDBNotification(n, *u, *fu)
	notif.ParentID = cDB.PostID.String()
	data, err := json.Marshal(notif)
	if err != nil {
		fail("failed to marshal JSON", err)
		return
	}
	if h.wsConnMann.IsConnected(u.Username) {
		h.wsConnMann.SendMessage(u.Username, data)
	}
}
