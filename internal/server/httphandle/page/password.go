package page

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"

	"github.com/Anacardo89/lenic/internal/middleware"
	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/server/httphandle/redirect"
	"github.com/Anacardo89/lenic/internal/session"
)

func (h *PageHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile(filepath.Join(h.homeDir, "templates/forgot-password.html"))
	if err != nil {
		h.log.Error("/forgot-password - Could not parse template", "error", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}

type RecoverPasswdPage struct {
	User  *models.User
	Token string
}

// /recover-password/{encoded_username}
func (h *PageHandler) RecoverPassword(w http.ResponseWriter, r *http.Request) {
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
	token := r.URL.Query().Get("token")
	if token == "" {
		fail("invalid token", err, true, http.StatusBadRequest, "invalid token")
		return
	}
	// DB operations
	dbUser, err := h.db.GetUserByUserName(r.Context(), username)
	if err != nil {
		fail("dberr: could not get user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	// Response
	u := models.FromDBUser(dbUser)
	page := RecoverPasswdPage{
		User:  u,
		Token: token,
	}
	t, err := template.ParseFiles(filepath.Join(h.homeDir, "templates/recover-password.html"))
	if err != nil {
		fail("could not parse template", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	t.Execute(w, page)
}

func (h *PageHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
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
	// Response
	t, err := template.ParseFiles(filepath.Join(h.homeDir, "templates/authorized/change-password.html"))
	if err != nil {
		fail("could not parse template", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	t.Execute(w, session)
}
