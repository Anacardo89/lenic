package pages

import (
	"html/template"
	"net/http"
	"time"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

type ErrorPage struct {
	ErrorMsg string
}

func Error(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/error ", r.RemoteAddr)
	cookieVal, err := r.Cookie("errormsg")
	if err != nil {
		logger.Error.Println("/error - Could not get error msg: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	logger.Info.Println("/error val: ", cookieVal.Value)
	if cookieVal.Expires.After(time.Now()) {
		cookieVal.MaxAge = -1
		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}
	errpg := ErrorPage{
		ErrorMsg: cookieVal.Value,
	}
	t, err := template.ParseFiles("templates/error.html")
	if err != nil {
		logger.Error.Println("/error - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, errpg)
}
