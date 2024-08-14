package actions

import (
	"database/sql"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/auth"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

func LoginPOST(w http.ResponseWriter, r *http.Request) {
	var err error
	u := presentation.User{
		UserName: r.FormValue("user_name"),
		UserPass: r.FormValue("user_password"),
	}
	if !isValidInput(u.UserName) || !isValidInput(u.UserPass) {
		RedirectToError(w, r, "Invalid character in form")
		return
	}
	err = db.Dbase.QueryRow(query.SelectUserByName, u.UserName).Scan(&u.Id, &u.UserName, &u.HashedPass, &u.Active)
	if err == sql.ErrNoRows {
		RedirectToError(w, r, "User does not exist")
		return
	}
	if u.Active == 0 {
		RedirectToError(w, r, "User is not active, check your mail")
		return
	}
	if !auth.CheckPasswordHash(u.UserPass, u.HashedPass) {
		RedirectToError(w, r, "Password does not match")
		return
	}
	usrSession := auth.CreateSession(w, r)
	auth.UpdateSession(usrSession.SessionId, u.Id)
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
