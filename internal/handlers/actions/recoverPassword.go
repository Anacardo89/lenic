package actions

import (
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

func RecoverPassword(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/recover-password ", r.RemoteAddr)
	// Parse Form
	err := r.ParseForm()
	if err != nil {
		logger.Error.Println("/action/recover-password - Could not parse Form: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userName := r.FormValue("user_name")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")

	if password != password2 {
		http.Error(w, "Password strings don't match", http.StatusBadRequest)
		return
	}

	hashed, err := auth.HashPassword(password)
	if err != nil {
		logger.Error.Println("/action/recover-password - Could not hash password: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = orm.Da.SetNewPassword(userName, hashed)
	if err != nil {
		logger.Error.Println("/action/recover-password - Could not set new password: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info.Println("OK - /action/recover-password ", r.RemoteAddr)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
