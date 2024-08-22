package pages

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

func Login(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/login ", r.RemoteAddr)
	body, err := os.ReadFile("templates/login.html")
	if err != nil {
		logger.Error.Println("/login - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}

func Register(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/register ", r.RemoteAddr)
	body, err := os.ReadFile("templates/register.html")
	if err != nil {
		logger.Error.Println("/register - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/forgot-password ", r.RemoteAddr)
	body, err := os.ReadFile("templates/forgot-password.html")
	if err != nil {
		logger.Error.Println("/forgot-password - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}

func RecoverPassword(w http.ResponseWriter, r *http.Request) {
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
	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Println("/recover-password - Could not get user: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t, err := template.ParseFiles("templates/recover-password.html")
	if err != nil {
		logger.Error.Println("/recover-password - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, dbuser)
}
