package actions

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	session := auth.ValidateSession(w, r)

	c := database.Comment{
		PostGUID:      postGUID,
		CommentAuthor: session.User.UserName,
		CommentText:   r.FormValue("comment_text"),
		Active:        1,
	}

	err := orm.Da.CreateComment(&c)
	if err != nil {
		logger.Error.Println(err)
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%s", postGUID), http.StatusSeeOther)
}

func EditComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["comment_id"])
	if err != nil {
		logger.Error.Println(err)
	}
	err = r.ParseForm()
	if err != nil {
		logger.Error.Println(err)
	}
	c := database.Comment{
		Id:          id,
		CommentText: r.FormValue("comment"),
	}
	if c.CommentText == "" {
		http.Error(w, "All form fields must be filled out", http.StatusBadRequest)
		return
	}

	orm.Da.UpdateCommentText(c.Id, c.CommentText)
	if err != nil {
		logger.Error.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}
