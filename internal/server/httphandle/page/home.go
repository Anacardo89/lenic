package page

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/Anacardo89/lenic/internal/server/httphandle/redirect"
	"github.com/Anacardo89/lenic/internal/session"
)

type HomePage struct {
	Session *session.Session
}

func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	feed := HomePage{
		Session: h.sm.ValidateSession(w, r),
	}
	t, err := template.ParseFiles("../frontend/templates/home.html")
	if err != nil {
		h.log.Error("/home - Could not parse template", "error", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, feed)
}

func (h *PageHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile("../frontend/templates/login.html")
	if err != nil {
		h.log.Error("/login - Could not parse template", "error", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}

func (h *PageHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile("../frontend/templates/register.html")
	if err != nil {
		h.log.Error("/register - Could not parse template", "error", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}
