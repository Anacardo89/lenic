package page

import (
	"encoding/base64"
	"html/template"
	"net/http"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/server/httphandle/redirect"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/gorilla/mux"
)

type FeedPage struct {
	Session *session.Session
	Posts   []*models.Post
}

func (h *PageHandler) Feed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("/user/%s/feed - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)

	dbUser, err := h.db.GetUserByUserName(h.ctx, userName)
	if err != nil {
		logger.Error.Printf("/user/%s/feed - Could not get user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	feed := FeedPage{}
	feed.Session = h.sessionStore.ValidateSession(w, r)
	dbPosts, err := h.db.GetFeed(h.ctx, dbUser.ID)
	if err != nil {
		logger.Error.Printf("/user/%s/feed - Could not get Posts: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	for _, p := range dbPosts {
		user, err := h.db.GetUserByID(h.ctx, p.AuthorID)
		if err != nil {
			logger.Error.Printf("/post/%s - Could not get Comment Author: %s\n", p.ID, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
		u := models.FromDBUser(user)
		post := models.FromDBPost(p, u)
		post.Content = template.HTML(post.RawContent)
		feed.Posts = append(feed.Posts, post)
	}
	t, err := template.ParseFiles("templates/authorized/feed.html")
	if err != nil {
		logger.Error.Printf("/user/%s/feed - Could not parse template: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, feed)
}
