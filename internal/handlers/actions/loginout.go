package actions

import (
	"database/sql"
	"net/http"
	"time"

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
		cookie := http.Cookie{Name: "errormsg",
			Value:    "Invalid character in form",
			Expires:  time.Now().Add(60 * time.Second),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/error", http.StatusMovedPermanently)
		return
	}
	err = db.Dbase.QueryRow(query.SelectUserByName, u.UserName).Scan(&u.Id, &u.UserName, &u.HashedPass, &u.Active)
	if err == sql.ErrNoRows {
		cookie := http.Cookie{Name: "errormsg",
			Value:    "User does not exist",
			Expires:  time.Now().Add(60 * time.Second),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/error", http.StatusMovedPermanently)
		return
	}
	if u.Active == 0 {
		cookie := http.Cookie{Name: "errormsg",
			Value:    "User is not active, check your mail",
			Expires:  time.Now().Add(60 * time.Second),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/error", http.StatusMovedPermanently)
		return
	}
	if !auth.CheckPasswordHash(u.UserPass, u.HashedPass) {
		cookie := http.Cookie{Name: "errormsg",
			Value:    "Password does not match",
			Expires:  time.Now().Add(60 * time.Second),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/error", http.StatusMovedPermanently)
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
