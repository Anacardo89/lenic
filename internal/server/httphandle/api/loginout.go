package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Anacardo89/lenic/internal/server/httphandle/redirect"
	"github.com/Anacardo89/lenic/pkg/crypto"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// /action/login
func (h *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
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
	var body LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fail("could not parse JSON from body", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// DB operations
	uDB, err := h.db.GetUserByUserName(h.ctx, body.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			fail("dberr: could not get user", err, true, http.StatusBadRequest, "invalid params")
			return
		} else {
			fail("dberr: could not get user", err, true, http.StatusInternalServerError, "internal error")
			return
		}
	}
	// Validation
	if !uDB.IsActive {
		fail("user is inactive", err, true, http.StatusUnauthorized, "inactive user")
		return
	}
	if err := crypto.ValidatePassword(uDB.PasswordHash, body.Password); err != nil {
		fail("wrong password", err, true, http.StatusUnauthorized, "wrong password")
		return
	}
	// Response
	if _, err := h.sm.CreateSession(w, r, uDB.ID); err != nil {
		fail("could not create session", err, true, http.StatusInternalServerError, "internal error")
	}
	w.WriteHeader(http.StatusOK)
}

// /action/logout
func (h *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.sm.DeleteSession(w, r)
	redirect.RedirIndex(w, r)
}
