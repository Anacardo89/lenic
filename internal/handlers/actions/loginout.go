package actions

import (
	"database/sql"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var err error
	userName := r.FormValue("user_name")
	userPass := r.FormValue("user_password")

	dbuser, err := orm.Da.GetUserByName(userName)
	if err == sql.ErrNoRows {
		RedirectToError(w, r, "User does not exist")
		return
	}
	logger.Info.Println(dbuser)
	u := mapper.User(dbuser)
	if u.Active != 1 {
		if u.Active == 2 {
			RedirectToError(w, r, "User is blocked, contact the admin")
			return
		}
		RedirectToError(w, r, "User is not active, check your mail")
		return
	}
	u.UserPass = userPass
	if !auth.CheckPasswordHash(u.UserPass, u.HashedPass) {
		RedirectToError(w, r, "Password does not match")
		return
	}
	usrSession := auth.CreateSession(w, r)
	auth.UpdateSession(usrSession.SessionId, u.Id)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := auth.SessionStore.Get(r, "tpsi25blog")
	if err != nil {
		logger.Error.Println(err)
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
