package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Anacardo89/lenic/internal/models"
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
	// DB operations
	notifsDB, err := h.db.GetUserNotifs(r.Context(), username, limit, offset)
	if err != nil {
		fail("dberr: could not get notifs", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// Response
	var notifs []*models.Notification
	for _, notif := range notifsDB {
		u := models.FromDBUserNotif(&notif.User)
		fromU := models.FromDBUserNotif(&notif.FromUser)
		n := models.FromDBNotification(&notif.Notification, *u, *fromU)
		notifs = append(notifs, n)
	}
	data, err := json.Marshal(notifs)
	if err != nil {
		fail("failed to marshal response body", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// PUT /action/user/{user_encoded}/notifications/{notif_id}/read
func (h *APIHandler) UpdateNotif(w http.ResponseWriter, r *http.Request) {
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
	nID, err := uuid.Parse(vars["notif_id"])
	if err != nil {
		fail("could not decode notif_id", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// DB operations
	err = h.db.UpdateNotificationRead(r.Context(), nID)
	if err != nil {
		fail("dberr: could not get notifs", err, true, http.StatusBadRequest, "invalid params")
		return
	}
}
