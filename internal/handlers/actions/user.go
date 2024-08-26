package actions

import (
	"encoding/base64"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

func FollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("POST /action/user/%s/follow %s\n", encoded, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/follow - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("POST /action/user/%s/follow %s %s\n", encoded, r.RemoteAddr, userName)

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/follow - Could not get user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	session := auth.ValidateSession(w, r)

	err = orm.Da.FollowUser(session.User.Id, dbuser.Id)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/follow - Could not follow: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("POST /action/user/%s/unfollow %s\n", encoded, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/unfollow - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("POST /action/user/%s/unfollow %s %s\n", encoded, r.RemoteAddr, userName)

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/unfollow - Could not get user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	session := auth.ValidateSession(w, r)

	err = orm.Da.UnfollowUser(session.User.Id, dbuser.Id)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/unfollow - Could not unfollow: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
