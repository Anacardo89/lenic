package redirect

import (
	"net/http"
	"time"
)

func RedirIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func RedirectToError(w http.ResponseWriter, r *http.Request, err string) {
	cookie := http.Cookie{Name: "errormsg",
		Value:    err,
		Expires:  time.Now().Add(60 * time.Second),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/error", http.StatusFound)
}
