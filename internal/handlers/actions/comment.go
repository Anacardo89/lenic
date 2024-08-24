package actions

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	logger.Info.Printf("POST /action/post/%s/comment %s\n", postGUID, r.RemoteAddr)
	session := auth.ValidateSession(w, r)

	c := database.Comment{
		PostGUID:  postGUID,
		AuthorId:  session.User.Id,
		Content:   r.FormValue("comment_text"),
		VoteCount: 0,
		Active:    1,
	}

	err := orm.Da.CreateComment(&c)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment - Could not create comment: %s\n", postGUID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%s", postGUID), http.StatusSeeOther)
}

func EditComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	id := vars["comment_id"]
	logger.Info.Printf("PUT /action/post/%s/comment/%s %s\n", postGUID, id, r.RemoteAddr)

	idint, err := strconv.Atoi(id)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment/%s - Could not convert id to string: %s\n", postGUID, id, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	err = r.ParseForm()
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment/%s - Could not parse form: %s\n", postGUID, id, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	c := database.Comment{
		Id:      idint,
		Content: r.FormValue("comment"),
	}
	if c.Content == "" {
		redirect.RedirectToError(w, r, "All form fields must be filled out")
		return
	}

	err = orm.Da.UpdateCommentText(c.Id, c.Content)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment/%s - Could not update comment: %s\n", postGUID, id, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	id := vars["comment_id"]
	logger.Info.Printf("DELETE /action/post/%s/comment/%s %s\n", postGUID, id, r.RemoteAddr)

	idint, err := strconv.Atoi(id)
	if err != nil {
		logger.Error.Printf("DELETE /action/post/%s/comment/%s - Could not convert id to string: %s\n", postGUID, id, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	err = orm.Da.DisableComment(idint)
	if err != nil {
		logger.Error.Printf("DELETE /action/post/%s/comment/%s - Could not update comment: %s\n", postGUID, id, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
