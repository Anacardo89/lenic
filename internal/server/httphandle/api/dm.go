package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/pkg/logger"
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
	// Auth
	session := h.sessionStore.ValidateSession(w, r)
	if !session.IsAuthenticated {
		fail("unauthorized", errors.New("unauthorized"), true, http.StatusUnauthorized, "unauthorized")
		return
	}
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
		fail("could not parse JSON from body", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	// DB operations
	c, users, err := h.db.GetConversationAndUsers(r.Context(), string(userBytes), msg.User)
	if err != nil {
		fail("dberr - could not get conversation", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	// Response
	var u, fromU *models.UserNotif
	if users[0].UserName == msg.User {
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
	// Auth
	session := h.sessionStore.ValidateSession(w, r)
	if !session.IsAuthenticated {
		fail("unauthorized", errors.New("unauthorized"), true, http.StatusUnauthorized, "unauthorized")
		return
	}
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
	// Auth
	session := h.sessionStore.ValidateSession(w, r)
	if !session.IsAuthenticated {
		fail("unauthorized", errors.New("unauthorized"), true, http.StatusUnauthorized, "unauthorized")
		return
	}
	// Input validation
	vars := mux.Vars(r)
	cID, err := uuid.Parse(vars["conversation_id"])
	if err != nil {
		fail("could not decode conversation", err, true, http.StatusBadRequest, "invalid conversation")
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]
	cIDstr := vars["conversation_id"]

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not decode user: %s\n", encoded, cIDstr, err)
		return
	}
	userName := string(bytes)

	cID, err := uuid.Parse(cIDstr)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not convert id to string: %s\n", encoded, cIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var msg JSON_DM
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not parse Json Data: %s\n", encoded, cIDstr, err)
		return
	}

	dbConvo, err := h.db.GetConversation(h.ctx, cID)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not get db conversation: %s\n", encoded, cIDstr, err)
		return
	}

	dbUser, err := h.db.GetUserByUserName(h.ctx, userName)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not get db user: %s\n", encoded, cIDstr, err)
		return
	}

	senderID := dbConvo.User1ID
	if dbConvo.User1ID != dbUser.ID {
		senderID = dbConvo.User2ID
	}

	m := &repo.DMessage{
		ConversationID: cID,
		SenderID:       senderID,
		Content:        msg.Msg,
	}

	_, err = h.db.CreateDM(h.ctx, m)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not get db user: %s\n", encoded, cIDstr, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// PUT /action/user/{user_encoded}/conversations/{conversation_id}/read
func (h *APIHandler) ReadConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]
	cIDstr := vars["conversation_id"]

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/conversations/%s/read - Could not decode user: %s\n", encoded, cIDstr, err)
		return
	}
	userName := string(bytes)

	dbUser, err := h.db.GetUserByUserName(h.ctx, userName)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/conversations/%s/read - Could get user: %s\n", encoded, cIDstr, err)
		return
	}

	cID, err := uuid.Parse(cIDstr)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/conversations/%s/read - Could not convert id to string: %s\n", encoded, cIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dms, err := h.db.GetDMsByConversation(h.ctx, cID, 1000, 0)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/conversations/%s/read - Could get DMs: %s\n", encoded, cIDstr, err)
		return
	}

	for _, dm := range dms {
		if dm.SenderID != dbUser.ID {
			err = h.db.UpdateDMRead(h.ctx, dm.ID)
			if err != nil {
				logger.Error.Printf("PUT /action/user/%s/conversations/%s/read - Could not update notif: %s\n", encoded, cIDstr, err)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}
