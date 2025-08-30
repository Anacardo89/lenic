package page

import (
	"database/sql"
	"encoding/base64"
	"html/template"
	"net/http"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/internal/server/httphandle/redirect"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/gorilla/mux"
)

type ProfilePage struct {
	Session *session.Session
	User    *models.User
	Posts   []*models.Post
	Follows string
}

func (h *PageHandler) UserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("/user/%s - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)

	dbUser, err := h.db.GetUserByUserName(h.ctx, userName)
	if err != nil {
		logger.Error.Printf("/user/%s - Could not get user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	u := models.FromDBUser(dbUser)
	session := h.sessionStore.ValidateSession(w, r)

	pp := ProfilePage{
		User:    u,
		Session: session,
	}

	dbFollow, err := h.db.GetUserFollows(h.ctx, session.User.ID, u.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error.Printf("/user/%s - Could not get follows: %s\n", encoded, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		} else {
			pp.Follows = models.StatusPending.String()
		}
	} else if dbFollow != nil {
		pp.Follows = dbFollow.FollowStatus
	} else {
		pp.Follows = models.StatusPending.String()
	}

	var dbPosts []*repo.Post
	if (session.User.ID == u.ID) || (dbFollow != nil && dbFollow.FollowStatus == models.StatusAccepted.String()) {
		dbPosts, err = h.db.GetUserPosts(h.ctx, u.ID)
		if err != nil {
			logger.Error.Printf("/user/%s - Could not get Posts: %s\n", encoded, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
	} else {
		dbPosts, err = h.db.GetUserPublicPosts(h.ctx, u.ID)
		if err != nil {
			logger.Error.Printf("/user/%s - Could not get Posts: %s\n", encoded, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
	}

	for _, p := range dbPosts {
		dbUser, err := h.db.GetUserByID(h.ctx, p.AuthorID)
		if err != nil {
			logger.Error.Printf("/post/%s - Could not get Comment Author: %s\n", p.ID, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
		u := models.FromDBUser(dbUser)
		post := models.FromDBPost(p, u)
		post.Content = template.HTML(post.RawContent)
		pp.Posts = append(pp.Posts, post)
	}

	t, err := template.ParseFiles("templates/authorized/user-profile.html")
	if err != nil {
		logger.Error.Printf("/user/%s - Could not parse template: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, pp)
}
