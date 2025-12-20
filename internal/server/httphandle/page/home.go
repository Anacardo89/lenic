package page

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Anacardo89/lenic/internal/server/httphandle/redirect"
)

func (h *PageHandler) Home(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile(filepath.Join(h.homeDir, "templates/home.html"))
	if err != nil {
		h.log.Error("/home - Could not read template", "error", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}

func (h *PageHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile(filepath.Join(h.homeDir, "templates/login.html"))
	if err != nil {
		h.log.Error("/login - Could not read template", "error", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}

func (h *PageHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile(filepath.Join(h.homeDir, "templates/register.html"))
	if err != nil {
		h.log.Error("/register - Could not read template", "error", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fmt.Fprint(w, string(body))
}
