package wshandle

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/repo"
)

// Endpoints:
//
// POST /action/post
// PUT /action/post/{post_id}
//
// ws - post_tag
func (h *WSHandler) HandlePostTag(msg Message, taggedUser string) {
	// Error Handling
	fail := func(logMsg string, e error) {
		h.log.Error(logMsg, "error", e,
			"message_type", msg.Type,
		)
	}
	//

	// Execution
	// Early return
	if taggedUser == msg.FromUserName {
		return
	}
	// DB operations
	uDB, err := h.db.GetUserByUserName(h.ctx, taggedUser)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	fuDB, err := h.db.GetUserByUserName(h.ctx, msg.FromUserName)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	noParent := ""
	n := &repo.Notification{
		UserID:     uDB.ID,
		FromUserID: fuDB.ID,
		NotifType:  msg.Type,
		NotifText:  msg.Msg,
		ResourceID: msg.ResourceID,
		ParentID:   &noParent,
	}
	if err := h.db.CreateNotification(h.ctx, n); err != nil {
		fail("dberr: could not create notification", err)
		return
	}
	// Response
	u := models.FromDBUserNotif(uDB)
	fromU := models.FromDBUserNotif(fuDB)
	notif := models.FromDBNotification(n, *u, *fromU)
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

// Endpoints:
//
// POST /action/post/{post_id}/comment
// PUT /action/post/{post_id}/comment/{comment_id}
//
// ws - comment_tag
func (h *WSHandler) HandleCommentTag(msg Message, taggedUser string) {
	// Error Handling
	fail := func(logMsg string, e error) {
		h.log.Error(logMsg, "error", e,
			"message_type", msg.Type,
		)
	}
	//

	// Execution
	// Early return
	if taggedUser == msg.FromUserName {
		return
	}
	cID, err := uuid.Parse(msg.ResourceID)
	if err != nil {
		fail("parsing comment uuid", err)
		return
	}
	// DB operations
	uDB, err := h.db.GetUserByUserName(h.ctx, taggedUser)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	fuDB, err := h.db.GetUserByUserName(h.ctx, msg.FromUserName)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	cDB, err := h.db.GetComment(h.ctx, cID)
	if err != nil {
		fail("dberr: could not get comment", err)
		return
	}
	n := &repo.Notification{
		UserID:     uDB.ID,
		FromUserID: fuDB.ID,
		NotifType:  msg.Type,
		NotifText:  msg.Msg,
		ResourceID: msg.ResourceID,
		ParentID:   &msg.ParentID,
	}
	if err := h.db.CreateNotification(h.ctx, n); err != nil {
		fail("dberr: could not create notification", err)
		return
	}
	// Response
	u := models.FromDBUserNotif(uDB)
	fromU := models.FromDBUserNotif(fuDB)
	notif := models.FromDBNotification(n, *u, *fromU)
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
