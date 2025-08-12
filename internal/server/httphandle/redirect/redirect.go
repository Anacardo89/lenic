package redirect

import (
	"net"
	"net/http"
	"time"

	"github.com/Anacardo89/lenic/pkg/logger"
)

func (h *RedirectHandler) RedirectNonSecure(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("Redirecting to HTTPS: ", r.RemoteAddr)
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = r.Host
	}
	redirectURL := "https://" + host + ":" + h.cfg.HTTPSPort + r.RequestURI
	logger.Info.Println(redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func (h *RedirectHandler) RedirIndex(w http.ResponseWriter, r *http.Request) {
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
