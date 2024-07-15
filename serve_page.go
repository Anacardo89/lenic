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

type ErrorPage struct {
	ErrorMsg string
}

func Index(w http.ResponseWriter, r *http.Request) {
	index := IndexPage{}
	index.Session = auth.ValidateSession(r)
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		logger.Error.Println(err)
	}
	t.Execute(w, index)
}

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile("templates/login.html")
	if err != nil {
		logger.Error.Println(err)
	}
	fmt.Fprint(w, string(body))
}

func Register(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile("templates/register.html")
	if err != nil {
		logger.Error.Println(err)
	}
	fmt.Fprint(w, string(body))
}

func Error(w http.ResponseWriter, r *http.Request) {
	cookieVal, err := r.Cookie("errormsg")
	if err != nil {
		logger.Error.Println(err)
	}
	errpg := ErrorPage{
		ErrorMsg: cookieVal.Value,
	}
	t, err := template.ParseFiles("templates/error.html")
	if err != nil {
		logger.Error.Println(err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "errormsg",
		MaxAge: -1,
	})
	t.Execute(w, errpg)
}
