package actions

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

type LoginRequest struct {
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/login ", r.RemoteAddr)
	var (
		err      error
		loginReq LoginRequest
	)

	err = json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		logger.Error.Println("/action/login - Could not decode JSON: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	defer r.Body.Close()

	logger.Info.Printf("/action/login %s %s\n", r.RemoteAddr, loginReq.UserName)

	dbuser, err := orm.Da.GetUserByName(loginReq.UserName)
	if err == sql.ErrNoRows {
		redirect.RedirectToError(w, r, "User does not exist")
		return
	}
	u := mapper.User(dbuser)
	if u.Active != 1 {
		if u.Active == 2 {
			redirect.RedirectToError(w, r, "User is blocked, contact the admin")
			return
		}
		redirect.RedirectToError(w, r, "User is not active, check your mail")
		return
	}
	u.Pass = loginReq.UserPassword
	if !auth.CheckPasswordHash(u.Pass, u.HashPass) {
		redirect.RedirectToError(w, r, "Password does not match")
		return
	}
	usrSession := auth.CreateSession(w, r)
	auth.UpdateSession(usrSession.SessionId, u.Id)

	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/logout ", r.RemoteAddr)
	session, err := auth.SessionStore.Get(r, "tpsi25blog")
	if err != nil {
		logger.Error.Println("/action/logout - Could not get session: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		logger.Error.Println("/action/logout - Could not save session: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
