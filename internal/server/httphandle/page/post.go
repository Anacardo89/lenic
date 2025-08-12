package page

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/server/httphandle/redirect"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/gorilla/mux"
)

type PostPage struct {
	Session *session.Session
	Post    *models.Post
}

func (h *PageHandler) NewPost(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/newPost ", r.RemoteAddr)
	postp := PostPage{
		Session: h.sessionStore.ValidateSession(w, r),
		Post:    &models.Post{},
	}
	t, err := template.ParseFiles("templates/authorized/newPost.html")
	if err != nil {
		logger.Error.Println("/newPost - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, postp)
}

func (h *PageHandler) Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]
	logger.Info.Printf("/post/%s %s\n", postID, r.RemoteAddr)
	pp := PostPage{}
	dbPost, err := h.db.GetPost(h.ctx, postID)
	if err != nil {
		logger.Error.Printf("/post/%s - Could not get Post: %s\n", postID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	dbUser, err := h.db.GetUserByID(h.ctx, dbPost.AuthorID)
	if err != nil {
		logger.Error.Printf("/post/%s - Could not get Author: %s\n", postID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	u := models.FromDBUser(dbUser)
	p := models.FromDBPost(dbPost, u)

	pp.Session = h.sessionStore.ValidateSession(w, r)
	pr, err := h.db.GetPostUserRating(h.ctx, dbPost.ID, pp.Session.User.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error.Printf("/post/%s - Could not get rating: %s\n", postID, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
	}
	if pr != nil {
		p.UserRating = pr.RatingValue
	} else {
		p.UserRating = 0
	}
	p.Content = template.HTML(p.RawContent)
	p.Comments = []*models.Comment{}

	dbComments, err := h.db.GetCommentsByPost(h.ctx, p.ID)
	if err != nil {
		logger.Error.Println(err)
	}
	for _, dbComment := range dbComments {
		dbUser, err := h.db.GetUserByID(h.ctx, dbComment.AuthorID)
		if err != nil {
			logger.Error.Printf("/post/%s - Could not get Comment Author: %s\n", postID, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}

		u := models.FromDBUser(dbUser)
		c := models.FromDBComment(dbComment, u)

		cr, err := h.db.GetCommentUserRating(h.ctx, dbComment.ID, pp.Session.User.ID)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.Error.Printf("/post/%s - Could not get comment rating: %s\n", postID, err)
				redirect.RedirectToError(w, r, err.Error())
				return
			}
		}

		if cr != nil {
			c.UserRating = cr.RatingValue
		} else {
			c.UserRating = 0
		}
		p.Comments = append(p.Comments, c)
	}
	pp.Post = p
	t, err := template.ParseFiles("templates/authorized/post.html")
	if err != nil {
		logger.Error.Printf("/post/%s - Could not parse template: %s\n", postID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, pp)
}
