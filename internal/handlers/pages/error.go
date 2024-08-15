package pages

import (
	"html/template"
	"net/http"
	"time"

	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

type ErrorPage struct {
	ErrorMsg string
}

func Error(w http.ResponseWriter, r *http.Request) {
	cookieVal, err := r.Cookie("errormsg")
	if err != nil {
		logger.Error.Println(err)
	}
	if cookieVal.Expires.After(time.Now()) {
		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}
	errpg := ErrorPage{
		ErrorMsg: cookieVal.Value,
	}
	t, err := template.ParseFiles("../templates/error.html")
	if err != nil {
		logger.Error.Println(err)
	}
	t.Execute(w, errpg)
}
