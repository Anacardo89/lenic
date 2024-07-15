package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog.git/auth"
	"github.com/Anacardo89/tpsi25_blog.git/db"
	"github.com/Anacardo89/tpsi25_blog.git/logger"
)

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	var err error
	u := auth.User{
		UserName: r.FormValue("user_name"),
		UserPass: r.FormValue("user_password"),
	}
	err = db.Dbase.QueryRow(db.SelectLogin, u.UserName).Scan(&u.Id, &u.UserName, &u.HashedPass)
	if err == sql.ErrNoRows {
		fmt.Fprintln(w, "User does not exist")
		return
	}
	if !auth.CheckPasswordHash(u.UserPass, u.HashedPass) {
		fmt.Fprintln(w, "Password does not match")
		return
	}
	usrSession := auth.CreateSession(w, r)
	db.UpdateSession(usrSession.Id, u.Id)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func LogoutPOST(w http.ResponseWriter, r *http.Request) {
	session, err := auth.SessionStore.Get(r, "tpsi25blog")
	if err != nil {
		logger.Error.Println(err)
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func ValidateSession(r *http.Request) auth.Session {
	usrSession := auth.Session{}
	session, err := auth.SessionStore.Get(r, "tpsi25blog")
	if err != nil {
		logger.Error.Println(err)
	}
	if sid, valid := session.Values["sid"]; valid {
		user := db.GetSessionUID(sid.(string))
		usrSession.User = auth.User{
			Id:        user.Id,
			UserName:  user.UserName,
			UserEmail: user.UserEmail,
		}
		db.UpdateSession(sid.(string), user.Id)
		usrSession.Id = sid.(string)
		usrSession.Authenticated = true
	} else {
		usrSession.Authenticated = false
	}
	return usrSession
}
