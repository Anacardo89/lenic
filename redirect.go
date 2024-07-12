package main

import (
	"net"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog.git/logger"
)

func RedirectNonSecure(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		logger.Error.Println(err)
		host = r.Host
	}
	redirectURL := "https://" + host + httpsServer.Addr + r.RequestURI
	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func RedirIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
