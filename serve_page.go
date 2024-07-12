package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/Anacardo89/tpsi25_blog.git/auth"
	"github.com/Anacardo89/tpsi25_blog.git/logger"
)

type IndexPage struct {
	Session auth.Session
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	index := IndexPage{}
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		logger.Error.Println(err)
	}
	t.Execute(w, index)
}

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile("templates/login.html")
	if err != nil {
		logger.Error.Println(err)
	}
	fmt.Fprint(w, string(body))
}

func ServeRegister(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile("templates/register.html")
	if err != nil {
		logger.Error.Println(err)
	}
	fmt.Fprint(w, string(body))
}
