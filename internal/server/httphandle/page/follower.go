package page

import (
	"encoding/base64"
	"html/template"
	"net/http"

	"github.com/Anacardo89/lenic/internal/handlers/data/orm"
	"github.com/Anacardo89/lenic/internal/handlers/redirect"
	"github.com/Anacardo89/lenic/internal/model/mapper"
	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/pkg/auth"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/gorilla/mux"
)

type FollowersPage struct {
	Session   models.Session
	User      models.User
	Followers []models.User
}

type FollowingPage struct {
	Session   models.Session
	User      models.User
	Following []models.User
}

func (h *PageHandler) UserFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("/user/%s/followers %s\n", encoded, r.RemoteAddr)

	followersp := FollowersPage{
		Session: auth.ValidateSession(w, r),
	}

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("/user/%s/followers - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("/user/%s/followers %s %s\n", encoded, r.RemoteAddr, userName)

	dbfollowed, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("/user/%s/followers - Could not get dbfollowed: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	followed := mapper.User(dbfollowed)
	followersp.User = *followed

	dbfollowers, err := orm.Da.GetFollowers(dbfollowed.Id)
	if err != nil {
		logger.Error.Printf("/user/%s/followers - Could not get dbfollowers: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	for _, dbfollower := range *dbfollowers {
		dbuser, err := orm.Da.GetUserByID(dbfollower.FollowerId)
		if err != nil {
			logger.Error.Printf("/user/%s/followers - Could not get dbuser: %s\n", encoded, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
		u := mapper.User(dbuser)
		followersp.Followers = append(followersp.Followers, *u)
	}

	t, err := template.ParseFiles("templates/authorized/followers.html")
	if err != nil {
		logger.Error.Printf("/user/%s/followers - Could not parse template: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, followersp)
}

func (h *PageHandler) UserFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("/user/%s/following %s\n", encoded, r.RemoteAddr)

	followingp := FollowingPage{
		Session: auth.ValidateSession(w, r),
	}

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("/user/%s/following - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("/user/%s/following %s %s\n", encoded, r.RemoteAddr, userName)

	dbfollower, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("/user/%s/following - Could not get dbfollower: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	follower := mapper.User(dbfollower)
	followingp.User = *follower

	dbfollowing, err := orm.Da.GetFollowing(dbfollower.Id)
	if err != nil {
		logger.Error.Printf("/user/%s/following - Could not get dbfollowing: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	for _, dbfollower := range *dbfollowing {
		dbuser, err := orm.Da.GetUserByID(dbfollower.FollowedId)
		if err != nil {
			logger.Error.Printf("/user/%s/following - Could not get dbuser: %s\n", encoded, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
		u := mapper.User(dbuser)
		followingp.Following = append(followingp.Following, *u)
	}

	t, err := template.ParseFiles("templates/authorized/following.html")
	if err != nil {
		logger.Error.Printf("/user/%s/following - Could not parse template: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, followingp)
}
