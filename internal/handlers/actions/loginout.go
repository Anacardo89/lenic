package actions

import (
	"database/sql"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var err error
	userName := r.FormValue("user_name")
	userPass := r.FormValue("user_password")
	logger.Info.Printf("/action/login %s %s\n", r.RemoteAddr, userName)

	dbuser, err := orm.Da.GetUserByName(userName)
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
	u.Pass = userPass
	if !auth.CheckPasswordHash(u.Pass, u.HashPass) {
		redirect.RedirectToError(w, r, "Password does not match")
		return
	}
	usrSession := auth.CreateSession(w, r)
	auth.UpdateSession(usrSession.SessionId, u.Id)
	userFeedPath := "/user/" + u.EncodedName + "/feed"
	http.Redirect(w, r, userFeedPath, http.StatusMovedPermanently)
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
	session.Save(r, w)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
