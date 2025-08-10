package server

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

type HomePage struct {
	Session presentation.Session
}

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/home ", r.RemoteAddr)
	feed := HomePage{}
	feed.Session = auth.ValidateSession(w, r)
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		logger.Error.Println("/home - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, feed)
}

func (s *Server) PageLogin(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/login ", r.RemoteAddr)
	body, err := os.ReadFile("templates/login.html")
	if err != nil {
		logger.Error.Println("/login - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}

func (s *Server) PageRegister(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/register ", r.RemoteAddr)
	body, err := os.ReadFile("templates/register.html")
	if err != nil {
		logger.Error.Println("/register - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}

func (s *Server) PageForgotPassword(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/forgot-password ", r.RemoteAddr)
	body, err := os.ReadFile("templates/forgot-password.html")
	if err != nil {
		logger.Error.Println("/forgot-password - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}

type RecoverPasswdPage struct {
	User  presentation.User
	Token string
}

func (s *Server) PageRecoverPassword(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/recover-password ", r.RemoteAddr)
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Println("/recover-password - Could not decode user: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("/recover-password %s %s", r.RemoteAddr, userName)

	token := r.URL.Query().Get("token")
	if token == "" {
		logger.Error.Println("/recover-password - No token", err)
		return
	}
	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Println("/recover-password - Could not get user: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	u := mapper.User(dbuser)
	page := RecoverPasswdPage{
		User:  *u,
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

func (s *Server) PageChangePassword(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/change-password ", r.RemoteAddr)
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Println("/change-password - Could not decode user: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("/change-password %s %s", r.RemoteAddr, userName)

	session := auth.ValidateSession(w, r)
	t, err := template.ParseFiles("templates/authorized/change-password.html")
	if err != nil {
		logger.Error.Println("/recover-password - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, session)
}
