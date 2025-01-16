package actions

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

// /action/recover-password
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
	token := r.FormValue("token")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Println("/action/recover-password - Could not get db user: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbToken, err := orm.Da.GetTokenByUserId(dbuser.Id)
	if err == sql.ErrNoRows {
		logger.Error.Println("/action/recover-password - No Token")
		http.Error(w, "No token", http.StatusBadRequest)
		return
	} else if err != nil {
		logger.Error.Println("/action/recover-password - Could not get db token: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if time.Now().After(dbToken.UpdatedAt.Add(time.Duration(1) * time.Hour)) {
		logger.Error.Println("/action/recover-password - Token Expired")
		http.Error(w, "Token Expired", http.StatusBadRequest)
		return
	}

	if dbToken.Token != token {
		logger.Error.Println("/action/recover-password - Token doesn't match")
		http.Error(w, "Token doesn't match", http.StatusBadRequest)
		return
	}

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

	err = orm.Da.DeleteTokenByUserId(dbuser.Id)
	if err != nil {
		logger.Error.Println("/action/recover-password - Could not delete token: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info.Println("OK - /action/recover-password ", r.RemoteAddr)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

// /action/change-password
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/change-password ", r.RemoteAddr)
	// Parse Form
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error.Println("/action/change-password - Could not parse form: ", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	userName := r.FormValue("user_name")
	logger.Debug.Println(userName)
	old_password := r.FormValue("old_password")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")

	dbUser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Println("/action/change-password - Could not get user: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !auth.CheckPasswordHash(old_password, dbUser.HashPass) {
		logger.Error.Println("/action/change-password - old password doesn't match, User: ", userName)
		http.Error(w, "old password doesn't match", http.StatusBadRequest)
		return
	}

	if password != password2 {
		http.Error(w, "Password strings don't match", http.StatusBadRequest)
		return
	}

	hashed, err := auth.HashPassword(password)
	if err != nil {
		logger.Error.Println("/action/change-password - Could not hash password: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = orm.Da.SetNewPassword(userName, hashed)
	if err != nil {
		logger.Error.Println("/action/change-password - Could not set new password: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info.Println("OK - /action/change-password ", r.RemoteAddr)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
