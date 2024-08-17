package actions

import (
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/auth"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

func RecoverPassword(w http.ResponseWriter, r *http.Request) {
	// Parse Form
	err := r.ParseForm()
	if err != nil {
		logger.Error.Println(err)
		return
	}
	userName := r.FormValue("user_name")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")

	if password != password2 {
		RedirectToError(w, r, "Password strings don't match")
		return
	}

	hashed, err := auth.HashPassword(password)
	if err != nil {
		logger.Error.Println(err)
		return
	}

	err = orm.Da.SetNewPassword(userName, hashed)
	if err != nil {
		logger.Error.Println(err)
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
