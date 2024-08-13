package pages

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Anacardo89/tpsi25_blog/auth"
	"github.com/Anacardo89/tpsi25_blog/internal/query"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
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
		CommentText: r.FormValue("comment_text"),
	}

	_, err := db.Dbase.Exec(query.InsertComment,
		postGUID,
		c.UserName,
		c.CommentText,
		1,
	)
	if err != nil {
		logger.Error.Println(err)
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%s", postGUID), http.StatusSeeOther)
}

func CommentPUT(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["comment_id"])
	if err != nil {
		logger.Error.Println(err)
	}
	err = r.ParseForm()
	if err != nil {
		logger.Error.Println(err)
	}
	c := Comment{
		Id:          id,
		CommentText: r.FormValue("comment"),
	}
	if c.CommentText == "" {
		http.Error(w, "All form fields must be filled out", http.StatusBadRequest)
		return
	}

	_, err = db.Dbase.Exec(query.UpdateComment,
		c.CommentText,
		c.Id,
	)
	if err != nil {
		logger.Error.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}
