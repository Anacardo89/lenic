package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GET /action/user/{user_encoded}/notifications
func (h *APIHandler) GetNotifs(w http.ResponseWriter, r *http.Request) {
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

	//
	//
	// Continue working here
	//
	//

	uDB, err := h.db.GetUserByUserName(h.ctx, username)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not get user: %s\n", encoded, err)
		return
	}
	u := models.FromDBUserNotif(uDB)

	dbNotifs, err := h.db.GetNotificationsByUser(h.ctx, u.ID, limit, offset)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not get notifs: %s\n", encoded, err)
		return
	}

	var notifs []*models.Notification
	for _, nDB := range dbNotifs {
		fromUserDB, err := h.db.GetUserByID(h.ctx, nDB.FromUserID)
		if err != nil {
			logger.Error.Printf("GET /action/user/%s/notifications - Could not get user: %s\n", encoded, err)
			return
		}
		fromU := models.FromDBUserNotif(fromUserDB)
		n := models.FromDBNotification(nDB, *u, *fromU)
		notifs = append(notifs, n)
	}

	data, err := json.Marshal(notifs)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not marshal notifs: %s\n", encoded, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// PUT /action/user/{user_encoded}/notifications/{notif_id}/read
func (h *APIHandler) UpdateNotif(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	nIDstr := vars["notif_id"]

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/notifications/%s/read - Could not decode user: %s\n", encoded, nIDstr, err)
		return
	}
	userName := string(bytes)

	nID, err := uuid.Parse(nIDstr)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/notifications/%s/read - Could not convert id to string: %s\n", encoded, nIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.UpdateNotificationRead(h.ctx, nID)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/notifications/%s/read - Could not update notif: %s\n", encoded, nIDstr, err)
		return
	}
}
