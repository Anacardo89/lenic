package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Anacardo89/tpsi25_blog.git/auth"
	"github.com/Anacardo89/tpsi25_blog.git/db"
	"github.com/gorilla/mux"
)

type Comment struct {
	Id          int
	UserName    string
	CommentText string
	Date        string
}

func CommentPOST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	session := auth.ValidateSession(r)

	c := Comment{
		UserName:    session.User.UserName,
		CommentText: r.FormValue("comment"),
	}

	_, err := db.Dbase.Exec(db.InsertComment,
		postGUID,
		c.UserName,
		c.CommentText,
	)
	if err != nil {
		log.Println(err.Error())
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%s", postGUID), http.StatusSeeOther)
}

func CommentPUT(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
	}
	id, err := strconv.Atoi(vars["comment_id"])
	postGUID := vars["post_guid"]
	c := Comment{
		Id:          id,
		CommentText: r.FormValue("edit_comment"),
	}
	if c.CommentText == "" {
		http.Error(w, "All form fields must be filled out", http.StatusBadRequest)
		return
	}

	_, err = db.Dbase.Exec(db.UpdateComment,
		c.CommentText,
		c.Id,
	)
	http.Redirect(w, r, fmt.Sprintf("/post/%s", postGUID), http.StatusSeeOther)
}
