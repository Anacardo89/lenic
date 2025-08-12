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

type FollowersPage struct {
	Session   *session.Session
	User      *models.User
	Followers []*models.User
}

type FollowingPage struct {
	Session   *session.Session
	User      *models.User
	Following []*models.User
}

func (h *PageHandler) UserFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]
	logger.Info.Printf("/user/%s/followers %s\n", encoded, r.RemoteAddr)

	fp := FollowersPage{
		Session: h.sessionStore.ValidateSession(w, r),
	}

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("/user/%s/followers - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("/user/%s/followers %s %s\n", encoded, r.RemoteAddr, userName)

	dbFollowed, err := h.db.GetUserByUserName(h.ctx, userName)
	if err != nil {
		logger.Error.Printf("/user/%s/followers - Could not get dbfollowed: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	followed := models.FromDBUser(dbFollowed)
	fp.User = followed

	dbFollowers, err := h.db.GetFollowers(h.ctx, dbFollowed.ID)
	if err != nil {
		logger.Error.Printf("/user/%s/followers - Could not get dbfollowers: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	for _, f := range dbFollowers {
		dbUser, err := h.db.GetUserByID(h.ctx, f.FollowerID)
		if err != nil {
			logger.Error.Printf("/user/%s/followers - Could not get dbuser: %s\n", encoded, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
		u := models.FromDBUser(dbUser)
		fp.Followers = append(fp.Followers, u)
	}

	t, err := template.ParseFiles("templates/authorized/followers.html")
	if err != nil {
		logger.Error.Printf("/user/%s/followers - Could not parse template: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, fp)
}

func (h *PageHandler) UserFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]
	logger.Info.Printf("/user/%s/following %s\n", encoded, r.RemoteAddr)

	fp := FollowingPage{
		Session: h.sessionStore.ValidateSession(w, r),
	}

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("/user/%s/following - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("/user/%s/following %s %s\n", encoded, r.RemoteAddr, userName)

	dbFollower, err := h.db.GetUserByUserName(h.ctx, userName)
	if err != nil {
		logger.Error.Printf("/user/%s/following - Could not get dbfollower: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	follower := models.FromDBUser(dbFollower)
	fp.User = follower

	dbFollowing, err := h.db.GetFollowing(h.ctx, dbFollower.ID)
	if err != nil {
		logger.Error.Printf("/user/%s/following - Could not get dbfollowing: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	for _, dbFollower := range dbFollowing {
		dbUser, err := h.db.GetUserByID(h.ctx, dbFollower.FollowedID)
		if err != nil {
			logger.Error.Printf("/user/%s/following - Could not get dbuser: %s\n", encoded, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
		u := models.FromDBUser(dbUser)
		fp.Following = append(fp.Following, u)
	}

	t, err := template.ParseFiles("templates/authorized/following.html")
	if err != nil {
		logger.Error.Printf("/user/%s/following - Could not parse template: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, fp)
}
