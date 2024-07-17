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
	t, err := template.ParseFiles("templates/newPost.html")
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
		&p.User,
		&p.RawContent,
		&p.Date,
	)
	if err != nil {
		logger.Error.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}
	p.Content = template.HTML(p.RawContent)
	comments, err := db.Dbase.Query(db.SelectComments, p.GUID)
	if err != nil {
		logger.Error.Println(err)
	}
	for comments.Next() {
		var c api.Comment
		comments.Scan(
			&c.Id,
			&c.UserName,
			&c.CommentText,
			&c.Date,
		)
		p.Comments = append(p.Comments, c)
	}
	t, err := template.ParseFiles("templates/post.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	t.Execute(w, p)

}

func ActivateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userName := vars["user_name"]
	_, err := db.Dbase.Exec(db.UpdateUserActive, userName)
	if err != nil {
		logger.Error.Println(err)
	}
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
