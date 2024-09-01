package actions

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("GET /action/search/user", r.RemoteAddr)
	queryParams := r.URL.Query()
	username := queryParams.Get("username")
	logger.Info.Printf("GET /action/search/user %s %s\n", r.RemoteAddr, username)

	dbusers, err := orm.Da.GetSearchUsers(username)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error.Printf("GET /action/search/user %s - Could not get users: %s\n", username, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
	}
	if dbusers == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	var users []presentation.UserNotif
	for _, dbuser := range *dbusers {
		u := mapper.UserNotif(&dbuser)
		users = append(users, *u)
	}

	data, err := json.Marshal(users)
	if err != nil {
		logger.Error.Printf("GET /action/search/user %s - Could not marshal users: %s\n", username, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	logger.Info.Printf("OK - GET /action/search/user %s %s\n", r.RemoteAddr, username)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func RequestFollowUser(w http.ResponseWriter, r *http.Request) {
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
	logger.Info.Printf("OK - POST /action/user/%s/follow %s\n", encoded, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("DELETE /action/user/%s/unfollow %s\n", encoded, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("DELETE /action/user/%s/unfollow - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("DELETE /action/user/%s/unfollow %s %s\n", encoded, r.RemoteAddr, userName)

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("DELETE /action/user/%s/unfollow - Could not get user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	queryParams := r.URL.Query()
	requesterName := queryParams.Get("requester")

	dbrequester, err := orm.Da.GetUserByName(requesterName)
	if err != nil {
		logger.Error.Printf("DELETE /action/user/%s/unfollow - Could not get dbrequester: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	err = orm.Da.UnfollowUser(dbrequester.Id, dbuser.Id)
	if err != nil {
		logger.Error.Printf("DELETE /action/user/%s/unfollow - Could not unfollow: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	dbnotif, err := orm.Da.GetFollowNotification(dbuser.Id, dbrequester.Id)
	if err != nil {
		logger.Error.Printf("DELETE /action/user/%s/unfollow - Could not get notif: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	err = orm.Da.DeleteNotificationByID(dbnotif.Id)
	if err != nil {
		logger.Error.Printf("DELETE /action/user/%s/unfollow - Could not delete notif: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	logger.Info.Printf("OK - DELETE /action/user/%s/unfollow %s\n", encoded, r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
}

func AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("PUT /action/user/%s/accept %s\n", encoded, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("PUT /action/user/%s/accept %s %s\n", encoded, r.RemoteAddr, userName)

	err = r.ParseForm()
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not parse form: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	requesterName := r.FormValue("requester")

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	dbrequester, err := orm.Da.GetUserByName(requesterName)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	err = orm.Da.AcceptFollow(dbrequester.Id, dbuser.Id)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not accept follow: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	dbnotif, err := orm.Da.GetFollowNotification(dbuser.Id, dbrequester.Id)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not get notif: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	err = orm.Da.DeleteNotificationByID(dbnotif.Id)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not delete notif: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	logger.Info.Printf("OK - PUT /action/user/%s/accept %s\n", encoded, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}
