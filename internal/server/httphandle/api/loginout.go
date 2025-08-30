package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/server/httphandle/redirect"
	"github.com/Anacardo89/lenic/pkg/crypto"
	"github.com/Anacardo89/lenic/pkg/logger"
)

type LoginRequest struct {
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_password"`
}

// /action/login
func (h *APIHandler) Login(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		loginReq LoginRequest
	)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error.Println("/action/login - Error reading body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &loginReq)
	if err != nil {
		logger.Error.Println("/action/login - Could not decode JSON: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	uDB, err := h.db.GetUserByUserName(h.ctx, loginReq.UserName)
	if err == sql.ErrNoRows {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}
	u := models.FromDBUser(uDB)
	if !u.IsActive {
		http.Error(w, "User is not active, check your mail", http.StatusBadRequest)
		return
	}
	u.Pass = loginReq.UserPassword
	if !crypto.ValidatePassword(u.PasswordHash, u.Pass) {
		http.Error(w, "Password does not match", http.StatusBadRequest)
		return
	}
	h.sessionStore.CreateSession(w, r, u.ID)

	w.WriteHeader(http.StatusOK)
}

// /action/logout
func (h *APIHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.sessionStore.DeleteSession(r)
	redirect.RedirIndex(w, r)
}
