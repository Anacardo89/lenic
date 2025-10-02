package api

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Anacardo89/lenic/internal/helpers"
	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/pkg/crypto"
)

// /action/register
func (h *APIHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
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
	// Hash password
	hashed, err := crypto.HashPassword(r.FormValue("password"))
	if err != nil {
		fail("could not hash password", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	var u = &models.User{
		Username:     r.FormValue("username"),
		Email:        r.FormValue("email"),
		Pass:         r.FormValue("password"),
		PasswordHash: hashed,
		ProfilePic:   "",
	}
	uDB := models.ToDBUser(u)
	// DB operations
	_, err = h.db.CreateUser(r.Context(), uDB)
	if err != nil {
		if err == repo.ErrUserExists {
			fail("dberr: user exists", err, true, http.StatusConflict, "username already exists")
			return
		}
		fail("dberr: failed to insert user", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Send account activation email
	mailSubject, mailBody := helpers.BuildActivateAccountMail("localhost", h.cfg.Port, u.Username)
	errs := h.mail.Send([]string{u.Email}, mailSubject, mailBody)
	if len(errs) != 0 {
		for _, err := range errs {
			fail("could not send account activation email", err, true, http.StatusInternalServerError, "internal error")
		}
		return
	}
	// Response
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

// /action/activate
func (h *APIHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
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
	err = h.db.SetUserActive(r.Context(), username)
	if err != nil {
		fail("dberr: could not activate user", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Response
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
