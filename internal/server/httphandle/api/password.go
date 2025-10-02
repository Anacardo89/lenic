package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/Anacardo89/lenic/internal/helpers"
	"github.com/Anacardo89/lenic/pkg/crypto"
)

// POST /action/forgot-password
func (h *APIHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
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
	// DB operations
	uDB, err := h.db.GetUserByEmail(h.ctx, r.FormValue("email"))
	if err == sql.ErrNoRows {
		fail("dberr: could not get user", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// Generate token
	token, err := h.tokenManager.GenerateToken(uDB.ID.String())
	if err != nil {
		fail("could not generate token", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Send recovery email
	mailSubject, mailBody := helpers.BuildPasswordRecoveryMail("localhost", h.cfg.Port, uDB.Username, token)
	errs := h.mail.Send([]string{uDB.Email}, mailSubject, mailBody)
	if len(errs) != 0 {
		for _, err := range errs {
			fail("could not send password recovery email", err, true, http.StatusInternalServerError, "internal error")
		}
		return
	}
	// Response
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

// /action/recover-password
func (h *APIHandler) RecoverPassword(w http.ResponseWriter, r *http.Request) {
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
	if r.FormValue("password") != r.FormValue("password2") {
		fail("passwords do not match", errors.New("passwords do not match"), true, http.StatusBadRequest, "passwords do not match")
		return
	}
	claims, err := h.tokenManager.ValidateToken(r.FormValue("token"))
	if err != nil {
		fail("invalid token", err, true, http.StatusUnauthorized, "invalid token")
		return
	}
	uID, err := uuid.Parse(claims.UserID)
	if err != nil {
		fail("could not decode user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	// Hash password
	hashed, err := crypto.HashPassword(r.FormValue("password"))
	if err != nil {
		fail("could not hash password", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// DB operations
	if err := h.db.SetNewPassword(h.ctx, uID, hashed); err != nil {
		fail("dberr: could not set password", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Response
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

// /action/change-password
func (h *APIHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
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
	if r.FormValue("password") != r.FormValue("password2") {
		fail("passwords do not match", errors.New("passwords do not match"), true, http.StatusBadRequest, "passwords do not match")
		return
	}
	// Get user from DB
	uDB, err := h.db.GetUserByUserName(r.Context(), r.FormValue("username"))
	if err != nil {
		fail("dberr: could not get user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	// Validate old password
	if !crypto.ValidatePassword(uDB.PasswordHash, r.FormValue("old_password")) {
		fail("wrong password", err, true, http.StatusUnauthorized, "wrong password")
		return
	}
	// Hash new password
	hashed, err := crypto.HashPassword(r.FormValue("password"))
	if err != nil {
		fail("could not hash password", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Set new password
	if err := h.db.SetNewPassword(h.ctx, uDB.ID, hashed); err != nil {
		fail("dberr: could not set password", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Response
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
