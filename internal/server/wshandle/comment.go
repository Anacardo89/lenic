package wshandle

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/repo"
)

func (h *WSHandler) handleCommentOnPost(msg Message) {
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
	// DB operations
	uDB, err := h.db.GetPostAuthorFromComment(h.ctx, cID)
	if err != nil {
		fail("dberr: could not get post author", err)
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
		ResourceID: cID.String(),
		ParentID:   msg.ParentID,
	}
	if err := h.db.CreateNotification(h.ctx, n); err != nil {
		fail("dberr: could not create notification", err)
		return
	}
	// Response
	u := models.FromDBUserNotif(uDB)
	fromU := models.FromDBUserNotif(fuDB)
	notif := models.FromDBNotification(n, *u, *fromU)
	notif.ParentID = msg.ParentID
	data, err := json.Marshal(notif)
	if err != nil {
		fail("failed to marshal JSON", err)
		return
	}
	if h.wsConnMann.IsConnected(u.Username) {
		h.wsConnMann.SendMessage(u.Username, data)
	}
}
