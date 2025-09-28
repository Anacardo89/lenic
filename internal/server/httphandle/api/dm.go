package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ConvoStarter struct {
	User string `json:"to_user"`
}

// POST /action/user/{user_encoded}/conversations
func (h *APIHandler) StartConversation(w http.ResponseWriter, r *http.Request) {
	// Error Handling
	fail := func(logMsg string, e error, writeError bool, status int, outMsg string) {
		h.log.Error(logMsg, "error", e,
			"status_code", status,
			"method", r.Method,
			"path", r.URL.Path,
			"client_ip", r.RemoteAddr,
		)
		if writeError {
			http.Error(w, outMsg, status)
		}
	}
	//

	// Execution
	// Input validation
	vars := mux.Vars(r)
	userBytes, err := base64.URLEncoding.DecodeString(vars["encoded_username"])
	if err != nil {
		fail("could not decode user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	var msg ConvoStarter
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		fail("could not parse JSON from body", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// DB operations
	c, users, err := h.db.GetConversationAndUsers(r.Context(), string(userBytes), msg.User)
	if err != nil {
		fail("dberr - could not get conversation", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// Response
	var u, fromU *models.UserNotif
	if users[0].Username == msg.User {
		fromU = models.FromDBUserNotif(users[0])
		u = models.FromDBUserNotif(users[1])
	} else {
		u = models.FromDBUserNotif(users[0])
		fromU = models.FromDBUserNotif(users[1])
	}
	convo := models.FromDBConversation(c, *u, *fromU, false)
	data, err := json.Marshal(convo)
	if err != nil {
		fail("could not marshal response body", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// GET /action/user/{user_encoded}/conversations
func (h *APIHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	// Error Handling
	fail := func(logMsg string, e error, writeError bool, status int, outMsg string) {
		h.log.Error(logMsg, "error", e,
			"status_code", status,
			"method", r.Method,
			"path", r.URL.Path,
			"client_ip", r.RemoteAddr,
		)
		if writeError {
			http.Error(w, outMsg, status)
		}
	}
	//

	// Execution
	// Input validation
	vars := mux.Vars(r)
	userBytes, err := base64.URLEncoding.DecodeString(vars["encoded_username"])
	if err != nil {
		fail("could not decode user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		fail("could not decode offset", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		fail("could not decode limit", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// DB operations
	dbUser, dbConvos, err := h.db.GetConversationsAndOwner(r.Context(), string(userBytes), limit, offset)
	if err != nil {
		fail("dberr - could not get user convos", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// Response
	u := models.FromDBUserNotif(dbUser)
	var convos []*models.Conversation
	for _, dbConvo := range dbConvos {
		isRead := true
		for _, dm := range dbConvo.Messages {
			if dm.SenderID != dbUser.ID && !dm.IsRead {
				isRead = false
				break
			}
		}
		c := models.FromDBConversationWithUser(dbConvo, *u, isRead)
		convos = append(convos, c)
	}
	data, err := json.Marshal(convos)
	if err != nil {
		fail("could not marshal response body", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// GET /action/user/{user_encoded}/conversations/{conversation_id}/dms
func (h *APIHandler) GetDMs(w http.ResponseWriter, r *http.Request) {
	// Error Handling
	fail := func(logMsg string, e error, writeError bool, status int, outMsg string) {
		h.log.Error(logMsg, "error", e,
			"status_code", status,
			"method", r.Method,
			"path", r.URL.Path,
			"client_ip", r.RemoteAddr,
		)
		if writeError {
			http.Error(w, outMsg, status)
		}
	}
	//

	// Execution
	// Input validation
	vars := mux.Vars(r)
	cID, err := uuid.Parse(vars["conversation_id"])
	if err != nil {
		fail("could not decode conversation", err, true, http.StatusBadRequest, "invalid conversation")
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		fail("could not decode offset", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		fail("could not decode limit", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// DB operations
	dbDMs, err := h.db.GetDMsByConversation(r.Context(), cID, limit, offset)
	if err != nil {
		fail("could not get DMs", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// Response
	var dms []*models.DMessage
	for _, dm := range dbDMs {
		d := repo.DMessage{
			ID:             dm.ID,
			ConversationID: dm.ConversationID,
			Content:        dm.Content,
			IsRead:         dm.IsRead,
			CreatedAt:      dm.CreatedAt,
		}
		sender := models.FromDBUserNotif(dm.Sender)
		dmOut := models.FromDBDMessage(&d, *sender)
		dms = append(dms, dmOut)
	}
	data, err := json.Marshal(dms)
	if err != nil {
		fail("could not marshal response body", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

type JSON_DM struct {
	Msg string `json:"text"`
}

// POST /action/user/{user_encoded}/conversations/{conversation_id}/dms
func (h *APIHandler) SendDM(w http.ResponseWriter, r *http.Request) {
	// Error Handling
	fail := func(logMsg string, e error, writeError bool, status int, outMsg string) {
		h.log.Error(logMsg, "error", e,
			"status_code", status,
			"method", r.Method,
			"path", r.URL.Path,
			"client_ip", r.RemoteAddr,
		)
		if writeError {
			http.Error(w, outMsg, status)
		}
	}
	//

	// Execution
	// Input validation
	vars := mux.Vars(r)
	bytes, err := base64.URLEncoding.DecodeString(vars["encoded_username"])
	if err != nil {
		fail("could not decode user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	username := string(bytes)
	cID, err := uuid.Parse(vars["conversation_id"])
	if err != nil {
		fail("could not decode conversation", err, true, http.StatusBadRequest, "invalid conversation")
		return
	}
	var msg JSON_DM
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		fail("could not decode body", err, true, http.StatusBadRequest, "invalid body")
		return
	}
	// DB operations
	cDB, uDB, err := h.db.GetConversationAndSender(r.Context(), cID, username)
	if err != nil {
		fail("dberr: could not get conversation", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	senderID := cDB.User1ID
	if cDB.User1ID != uDB.ID {
		senderID = cDB.User2ID
	}
	m := &repo.DMessage{
		ConversationID: cID,
		SenderID:       senderID,
		Content:        msg.Msg,
	}
	go func() {
		if _, err = h.db.CreateDM(context.Background(), m); err != nil {
			fail("dberr: could not create dm", err, true, http.StatusInternalServerError, "internal error")
			return
		}
	}()
	// Response
	w.WriteHeader(http.StatusOK)
}

// PUT /action/user/{user_encoded}/conversations/{conversation_id}/read
func (h *APIHandler) ReadConversation(w http.ResponseWriter, r *http.Request) {
	// Error Handling
	fail := func(logMsg string, e error, writeError bool, status int, outMsg string) {
		h.log.Error(logMsg, "error", e,
			"status_code", status,
			"method", r.Method,
			"path", r.URL.Path,
			"client_ip", r.RemoteAddr,
		)
		if writeError {
			http.Error(w, outMsg, status)
		}
	}
	//

	// Execution
	// Input validation
	vars := mux.Vars(r)
	bytes, err := base64.URLEncoding.DecodeString(vars["encoded_username"])
	if err != nil {
		fail("could not decode user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	username := string(bytes)
	cID, err := uuid.Parse(vars["conversation_id"])
	if err != nil {
		fail("could not decode conversation", err, true, http.StatusBadRequest, "invalid conversation")
		return
	}
	// DB operations
	if err := h.db.ReadAllReceivedDMsInConvo(h.ctx, cID, username); err != nil {
		fail("dberr: could not decode read DMs", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// Response
	w.WriteHeader(http.StatusOK)
}
