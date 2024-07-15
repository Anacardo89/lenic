package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/Anacardo89/tpsi25_blog.git/api"
	"github.com/Anacardo89/tpsi25_blog.git/auth"
	"github.com/Anacardo89/tpsi25_blog.git/db"
	"github.com/Anacardo89/tpsi25_blog.git/logger"
	"github.com/gorilla/mux"
)

type IndexPage struct {
	Posts   []api.PostPage
	Session auth.Session
}

type ErrorPage struct {
	ErrorMsg string
}

func Index(w http.ResponseWriter, r *http.Request) {
	index := IndexPage{}
	index.Session = auth.ValidateSession(r)
	rows, err := db.Dbase.Query(db.SelectPosts)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		thisPost := api.PostPage{}
		err := rows.Scan(
			&thisPost.GUID,
			&thisPost.Title,
			&thisPost.User,
			&thisPost.RawContent,
			&thisPost.Date,
		)
		if err != nil {
			logger.Error.Println(err)
			return
		}
		thisPost.Content = template.HTML(thisPost.RawContent)
		index.Posts = append(index.Posts, thisPost)
	}
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

func NewPost(w http.ResponseWriter, r *http.Request) {
	postpg := api.PostPage{
		Session: auth.ValidateSession(r),
	}
	t, err := template.ParseFiles("templates/post.html")
	if err != nil {
		logger.Error.Println(err)
	}
	t.Execute(w, postpg)
}

func Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p := api.PostPage{
		Session: auth.ValidateSession(r),
		GUID:    vars["post_guid"],
	}
	err := db.Dbase.QueryRow(db.SelectPostByGUID, p.GUID).Scan(
		&p.Title,
		&p.RawContent,
		&p.Date,
	)
	if err != nil {
		logger.Error.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}
	p.Content = template.HTML(p.RawContent)
	t, err := template.ParseFiles("templates/posts.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	t.Execute(w, p)

}
