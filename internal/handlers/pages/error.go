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
	logger.Info.Println("/error ", r.RemoteAddr)

	queryParams := r.URL.Query()
	msg := queryParams.Get("message")

	errpg := ErrorPage{
		ErrorMsg: msg,
	}
	t, err := template.ParseFiles("templates/error.html")
	if err != nil {
		logger.Error.Println("/error - Could not parse template: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.Execute(w, errpg)
}
