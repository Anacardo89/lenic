package api

import (
	"database/sql"
	"net/http"

	"github.com/Anacardo89/lenic/internal/helpers"
	"github.com/Anacardo89/lenic/pkg/crypto"
	"github.com/Anacardo89/lenic/pkg/logger"
)

// POST /action/forgot-password
func (h *APIHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	// Parse Form
	err := r.ParseForm()
	if err != nil {
		logger.Error.Println("POST /action/forgot-password - Could not parse Form: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mail := r.FormValue("user_email")
	// Get user from DB
	dbuser, err := h.db.GetUserByEmail(h.ctx, mail)
	if err == sql.ErrNoRows {
		http.Error(w, "No user with that email", http.StatusBadRequest)
		return
	}

	token, err := h.tokenManager.GenerateToken(dbuser.ID)
	if err != nil {
		logger.Error.Println("POST /action/forgot-password - Could not generate token: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mailSubject, mailBody := helpers.BuildPasswordRecoveryMail(h.mail.Host, string(h.mail.Port), dbuser.UserName, token)
	errs := h.mail.Send([]string{dbuser.Email}, mailSubject, mailBody)
	if len(errs) != 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

// /action/recover-password
func (h *APIHandler) RecoverPassword(w http.ResponseWriter, r *http.Request) {
	// Parse Form
	err := r.ParseForm()
	if err != nil {
		logger.Error.Println("/action/recover-password - Could not parse Form: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token := r.FormValue("token")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")

	claims, err := h.tokenManager.ValidateToken(token)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	dbuser, err := h.db.GetUserByID(h.ctx, claims.UserID)
	if err != nil {
		logger.Error.Println("/action/recover-password - Could not get db user: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if password != password2 {
		http.Error(w, "Password strings don't match", http.StatusBadRequest)
		return
	}

	hashed, err := crypto.HashPassword(password)
	if err != nil {
		logger.Error.Println("/action/recover-password - Could not hash password: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.SetNewPassword(h.ctx, dbuser.UserName, hashed)
	if err != nil {
		logger.Error.Println("/action/recover-password - Could not set new password: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

// /action/change-password
func (h *APIHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Parse Form
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error.Println("/action/change-password - Could not parse form: ", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	username := r.FormValue("user_name")
	old_password := r.FormValue("old_password")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")

	dbUser, err := h.db.GetUserByUserName(h.ctx, username)
	if err != nil {
		logger.Error.Println("/action/change-password - Could not get user: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !crypto.ValidatePassword(dbUser.PasswordHash, old_password) {
		logger.Error.Println("/action/change-password - old password doesn't match, User: ", username)
		http.Error(w, "old password doesn't match", http.StatusBadRequest)
		return
	}

	if password != password2 {
		http.Error(w, "Password strings don't match", http.StatusBadRequest)
		return
	}

	hashed, err := crypto.HashPassword(password)
	if err != nil {
		logger.Error.Println("/action/change-password - Could not hash password: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.SetNewPassword(h.ctx, username, hashed)
	if err != nil {
		logger.Error.Println("/action/change-password - Could not set new password: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
