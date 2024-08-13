package pages

import (
	"html/template"
	"net/http"

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
	errpg := ErrorPage{
		ErrorMsg: cookieVal.Value,
	}
	t, err := template.ParseFiles("../templates/error.html")
	if err != nil {
		logger.Error.Println(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "errormsg",
		MaxAge: -1,
	})
	t.Execute(w, errpg)
}
