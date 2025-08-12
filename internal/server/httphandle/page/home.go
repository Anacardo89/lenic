package page

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/Anacardo89/lenic/internal/server/httphandle/redirect"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/Anacardo89/lenic/pkg/logger"
)

type HomePage struct {
	Session *session.Session
}

func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/home ", r.RemoteAddr)
	feed := HomePage{}
	feed.Session = h.sessionStore.ValidateSession(w, r)
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		logger.Error.Println("/home - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, feed)
}

func (h *PageHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/login ", r.RemoteAddr)
	body, err := os.ReadFile("templates/login.html")
	if err != nil {
		logger.Error.Println("/login - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}

func (h *PageHandler) Register(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/register ", r.RemoteAddr)
	body, err := os.ReadFile("templates/register.html")
	if err != nil {
		logger.Error.Println("/register - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}
