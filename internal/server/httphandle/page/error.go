package page

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type ErrorPage struct {
	ErrorMsg string
}

func (h *PageHandler) Error(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	msg := queryParams.Get("message")

	errpg := ErrorPage{
		ErrorMsg: msg,
	}
	t, err := template.ParseFiles(filepath.Join(h.homeDir, "templates/error.html"))
	if err != nil {
		h.log.Error("/error - Could not parse template", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.Execute(w, errpg)
}
