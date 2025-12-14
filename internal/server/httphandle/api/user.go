package api

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Anacardo89/lenic/internal/middleware"
	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/gorilla/mux"
)

// GET /action/search/user
func (h *APIHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
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
	// DB operations
	usersDB, err := h.db.SearchUsersByUserName(r.Context(), r.URL.Query().Get("username"))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			fail("dberr: could not get users", err, true, http.StatusInternalServerError, "internal error")
			return
		}
	}
	// Response
	var users []models.UserNotif
	for _, uDB := range usersDB {
		u := models.FromDBUserNotif(uDB)
		users = append(users, *u)
	}
	data, err := json.Marshal(users)
	if err != nil {
		fail("could not marshal response body", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// POST /action/user/{user_encoded}/follow
func (h *APIHandler) RequestFollowUser(w http.ResponseWriter, r *http.Request) {
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
	// Get session
	session, ok := r.Context().Value(middleware.CtxKeySession).(*session.Session)
	if !ok {
		fail("session type mismatch", errors.New("session type mismatch"), true, http.StatusUnauthorized, "invalid session")
		return
	}
	// Input validation
	vars := mux.Vars(r)
	bytes, err := base64.URLEncoding.DecodeString(vars["encoded_username"])
	if err != nil {
		fail("could not decode user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	username := string(bytes)
	// DB operations
	if err := h.db.FollowUser(r.Context(), session.User.ID, username); err != nil {
		fail("dberr: could not follow user", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Response
	w.WriteHeader(http.StatusOK)
}

// DELETE /action/user/{user_encoded}/unfollow
func (h *APIHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
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
	// DB operations
	if err := h.db.UnfollowUser(r.Context(), r.URL.Query().Get("requester"), username); err != nil {
		fail("dberr: could not unfollow user", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	if err := h.db.DeleteFollowNotification(r.Context(), username, r.URL.Query().Get("requester")); err != nil {
		fail("dberr: could not delete follow notification", err, false, http.StatusInternalServerError, "internal error")
	}
	// Response
	w.WriteHeader(http.StatusOK)
}

// PUT /action/user/{user_encoded}/accept
func (h *APIHandler) AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
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
	if err := r.ParseForm(); err != nil {
		fail("could not parse form", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	vars := mux.Vars(r)
	bytes, err := base64.URLEncoding.DecodeString(vars["encoded_username"])
	if err != nil {
		fail("could not decode user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	username := string(bytes)
	// DB operations
	if err := h.db.AcceptFollow(h.ctx, r.FormValue("requester"), username); err != nil {
		fail("dberr: could not accept follow", err, false, http.StatusBadRequest, "invalid params")
		return
	}
	if err := h.db.DeleteFollowNotification(r.Context(), username, r.FormValue("requester")); err != nil {
		fail("dberr: could not delete follow notification", err, false, http.StatusInternalServerError, "internal error")
	}
	// Response
	w.WriteHeader(http.StatusOK)
}
