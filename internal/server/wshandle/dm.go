package wshandle

import (
	"encoding/json"

	"github.com/Anacardo89/lenic/internal/models"
)

// ws - dm
func (h *WSHandler) handleDM(msg Message) {
	// Error Handling
	fail := func(logMsg string, e error) {
		h.log.Error(logMsg, "error", e,
			"message_type", msg.Type,
		)
	}
	//

	// Execution
	// Early return
	if msg.ResourceID == msg.FromUserName {
		return
	}
	// DB ooperations
	uDB, err := h.db.GetUserByUserName(h.ctx, msg.ResourceID)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	fuDB, err := h.db.GetUserByUserName(h.ctx, msg.FromUserName)
	if err != nil {
		fail("dberr: could not get user", err)
		return
	}
	dbConvo, err := h.db.GetConversationByUsers(h.ctx, uDB.ID, fuDB.ID)
	if err != nil {
		fail("dberr: could not get conversation", err)
		return
	}
	if err := h.db.UpdateConversation(h.ctx, dbConvo.ID); err != nil {
		fail("dberr: could not update conversation", err)
		return
	}
	// Response
	u := models.FromDBUserNotif(uDB)
	fu := models.FromDBUserNotif(fuDB)
	n := &models.Notification{
		User:       *u,
		FromUser:   *fu,
		NotifType:  models.NotifType(msg.Type),
		NotifText:  msg.Msg,
		ResourceID: dbConvo.ID.String(),
		ParentID:   "",
		IsRead:     false,
	}
	data, err := json.Marshal(n)
	if err != nil {
		fail("failed to marshal JSON", err)
		return
	}
	if h.wsConnMann.IsConnected(u.Username) {
		h.wsConnMann.SendMessage(u.Username, data)
	}
}
