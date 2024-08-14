package actions

import (
	"net/http"
	"time"
)

func RedirectToError(w http.ResponseWriter, r *http.Request, err string) {
	cookie := http.Cookie{Name: "errormsg",
		Value:    err,
		Expires:  time.Now().Add(60 * time.Second),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/error", http.StatusMovedPermanently)
}
