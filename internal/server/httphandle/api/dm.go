package api

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type JSON_Convo struct {
	User string `json:"to_user"`
}

// POST /action/user/{user_encoded}/conversations
func (h *APIHandler) StartConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]
	logger.Info.Printf("POST /action/user/%s/conversations %s\n", encoded, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations - Could not decode user: %s\n", encoded, err)
		return
	}
	userName := string(bytes)
	logger.Info.Printf("POST /action/user/%s/conversations %s %s\n", encoded, r.RemoteAddr, userName)

	var msg JSON_Convo
	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations - Could not parse Json Data: %s\n", encoded, err)
		return
	}

	dbUser, err := h.db.GetUserByUserName(h.ctx, userName)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations - Could not get user: %s\n", encoded, err)
		return
	}
	u := models.FromDBUserNotif(dbUser)

	dbFromUser, err := h.db.GetUserByUserName(h.ctx, msg.User)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations - Could not get from user: %s\n", encoded, err)
		return
	}
	fromU := models.FromDBUserNotif(dbFromUser)

	exists := true
	var dbConvo *db.Conversation
	dbConvo, err = h.db.GetConversationByUsers(h.ctx, u.ID, fromU.ID)
	if err == sql.ErrNoRows {
		exists = false
	} else if err != nil {
		logger.Error.Println("Could not get conversation: ", err)
		return
	}
	if exists {
		convo := &db.Conversation{
			User1ID: u.ID,
			User2ID: fromU.ID,
		}
		convoID, err := h.db.CreateConversation(h.ctx, convo)
		if err != nil {
			logger.Error.Println("Could not create conversation: ", err)
			return
		}
		dbConvo, err = h.db.GetConversation(h.ctx, convoID)
		if err != nil {
			logger.Error.Printf("POST /action/user/%s/conversations - Could not get conversation: %s\n", encoded, err)
			return
		}
	}
	convo := models.FromDBConversation(dbConvo, *u, *fromU, false)

	data, err := json.Marshal(convo)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations - Could not marshal conversations: %s\n", encoded, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// GET /action/user/{user_encoded}/conversations
func (h *APIHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]
	logger.Info.Printf("GET /action/user/%s/conversations %s\n", encoded, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations - Could not decode user: %s\n", encoded, err)
		return
	}
	userName := string(bytes)
	logger.Info.Printf("GET /action/user/%s/conversations %s %s\n", encoded, r.RemoteAddr, userName)

	dbUser, err := h.db.GetUserByUserName(h.ctx, userName)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations - Could not get user: %s\n", encoded, err)
		return
	}
	u := models.FromDBUserNotif(dbUser)

	queryParams := r.URL.Query()
	offsetStr := queryParams.Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations - Could not parse offset to int: %s\n", encoded, err)
		return
	}

	limitStr := queryParams.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations - Could not parse limit to int: %s\n", encoded, err)
		return
	}

	dbConvos, err := h.db.GetConversationsByUser(h.ctx, dbUser.ID, limit, offset)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations - Could not get conversations: %s\n", encoded, err)
		return
	}

	var convos []*models.Conversation
	for _, dbConvo := range dbConvos {
		fromUserID := dbConvo.User1ID
		if dbConvo.User1ID == dbUser.ID {
			fromUserID = dbConvo.User2ID
		}
		dbFromUser, err := h.db.GetUserByID(h.ctx, fromUserID)
		if err != nil {
			logger.Error.Printf("GET /action/user/%s/conversations - Could not get user: %s\n", encoded, err)
			return
		}
		fromU := models.FromDBUserNotif(dbFromUser)
		dms, err := h.db.GetDMsByConversation(h.ctx, dbConvo.ID, 1000, 0)
		if err != nil {
			logger.Error.Printf("GET /action/user/%s/conversations - Could not get dms: %s\n", encoded, err)
			return
		}
		isRead := true
		for _, dm := range dms {
			if dm.SenderID != dbUser.ID && !dm.IsRead {
				isRead = false
				break
			}
		}
		c := models.FromDBConversation(dbConvo, *u, *fromU, isRead)
		convos = append(convos, c)
	}

	data, err := json.Marshal(convos)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations - Could not marshal conversations: %s\n", encoded, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// GET /action/user/{user_encoded}/conversations/{conversation_id}/dms
func (h *APIHandler) GetDMs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]
	cIDstr := vars["conversation_id"]
	logger.Info.Printf("GET /action/user/%s/conversations/%s/dms %s\n", encoded, cIDstr, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not decode user: %s\n", encoded, cIDstr, err)
		return
	}
	userName := string(bytes)
	logger.Info.Printf("GET /action/user/%s/conversations/%s/dms %s %s\n", encoded, cIDstr, r.RemoteAddr, userName)

	cID, err := uuid.Parse(cIDstr)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not convert id to string: %s\n", encoded, cIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	offsetStr := queryParams.Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not parse offset to int: %s\n", encoded, cIDstr, err)
		return
	}

	limitStr := queryParams.Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not parse limit to int: %s\n", encoded, cIDstr, err)
		return
	}

	dbConvo, err := h.db.GetConversation(h.ctx, cID)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not get conversation: %s\n", encoded, cIDstr, err)
		return
	}

	var dms []*models.DMessage
	dbDMs, err := h.db.GetDMsByConversation(h.ctx, dbConvo.ID, limit, offset)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not get DMs: %s\n", encoded, cIDstr, err)
		return
	}
	for _, dm := range dbDMs {
		senderDB, err := h.db.GetUserByID(h.ctx, dm.SenderID)
		if err != nil {
			logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not get sender: %s\n", encoded, cIDstr, err)
			return
		}
		sender := models.FromDBUserNotif(senderDB)
		dm := models.FromDBDMessage(dm, *sender)
		dms = append(dms, dm)
	}

	data, err := json.Marshal(dms)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/conversations/%s/dms - Could not marshal dms: %s\n", encoded, cIDstr, err)
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
	logger.Info.Printf("POST /action/user/%s/conversations/%s/dms %s\n", encoded, cIDstr, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not decode user: %s\n", encoded, cIDstr, err)
		return
	}
	userName := string(bytes)
	logger.Info.Printf("POST /action/user/%s/conversations/%s/dms %s %s\n", encoded, cIDstr, r.RemoteAddr, userName)

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

	m := &db.DMessage{
		ConversationID: cID,
		SenderID:       senderID,
		Content:        msg.Msg,
	}

	_, err = h.db.CreateDM(h.ctx, m)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/conversations/%s/dms - Could not get db user: %s\n", encoded, cIDstr, err)
		return
	}

	logger.Info.Printf("OK - POST /action/user/%s/conversations/%s/dms\n", encoded, cIDstr)
	w.WriteHeader(http.StatusOK)
}

// PUT /action/user/{user_encoded}/conversations/{conversation_id}/read
func (h *APIHandler) ReadConversation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]
	cIDstr := vars["conversation_id"]
	logger.Info.Printf("PUT /action/user/%s/conversations/%s/read %s\n", encoded, cIDstr, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/conversations/%s/read - Could not decode user: %s\n", encoded, cIDstr, err)
		return
	}
	userName := string(bytes)
	logger.Info.Printf("PUT /action/user/%s/conversations/%s/read %s %s\n", encoded, cIDstr, r.RemoteAddr, userName)

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

	logger.Info.Printf("OK - PUT /action/user/%s/conversations/%s/read\n", encoded, cIDstr)
	w.WriteHeader(http.StatusOK)
}
