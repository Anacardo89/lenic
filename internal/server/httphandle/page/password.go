package page

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/server/httphandle/redirect"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/gorilla/mux"
)

func (h *PageHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile("templates/forgot-password.html")
	if err != nil {
		logger.Error.Println("/forgot-password - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}

type RecoverPasswdPage struct {
	User  *models.User
	Token string
}

func (h *PageHandler) RecoverPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]
	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Println("/recover-password - Could not decode user: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)

	token := r.URL.Query().Get("token")
	if token == "" {
		logger.Error.Println("/recover-password - No token", err)
		return
	}
	dbUser, err := h.db.GetUserByUserName(h.ctx, userName)
	if err != nil {
		logger.Error.Println("/recover-password - Could not get user: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	u := models.FromDBUser(dbUser)
	page := RecoverPasswdPage{
		User:  u,
		Token: token,
	}
	t, err := template.ParseFiles("templates/recover-password.html")
	if err != nil {
		logger.Error.Println("/recover-password - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, page)
}

func (h *PageHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]
	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Println("/change-password - Could not decode user: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)

	session := h.sessionStore.ValidateSession(w, r)
	t, err := template.ParseFiles("templates/authorized/change-password.html")
	if err != nil {
		logger.Error.Println("/recover-password - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, session)
}
