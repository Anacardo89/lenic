package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/Anacardo89/tpsi25_blog.git/auth"
	"github.com/Anacardo89/tpsi25_blog.git/logger"
)

func isValidInput(input string) bool {
	if strings.Contains(input, ";") {
		return false
	}
	return true
}

func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	var userReg *auth.User
	userReg.UserName = r.FormValue("user_name")
	userReg.UserEmail = r.FormValue("user_email")
	userReg.UserPass = r.FormValue("user_password")
	pass2 := r.FormValue("user_password2")
	if userReg.UserPass != pass2 {
		return
	}
	if !isValidInput(userReg.UserName) || !isValidInput(userReg.UserEmail) || !isValidInput(userReg.UserPass) {
		return
	}

	// TODO
	// check repeated user_name and user_email

	userReg.HashedPass, err = auth.HashPassword(userReg.UserPass)
	if err != nil {
		return
	}

	// TODO
	// rework db package similar to kanboards
	_, err = database.Exec("INSERT INTO users SET user_name=?, user_guid=?, user_email=?, user_password=?", name, guid, email, password)
	if err != nil {
		logger.Error.Println(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	} else {
		Email := RegistrationData{Email: email, Message: ""}
		message, err := template.New("registrationEmail").Parse(WelcomeEmail)
		if err != nil {
			logger.Error.Println(err.Error())
			return
		}
		var mbuf bytes.Buffer
		err = message.Execute(&mbuf, Email)
		if err != nil {
			logger.Error.Println(err.Error())
			return
		}
		Email.Message = mbuf.String()
		data, err := json.Marshal(Email)
		if err != nil {
			logger.Error.Println(err.Error())
			return
		}
		err = MQSendMail(data)
		if err != nil {
			logger.Error.Println(err.Error())
			return
		}

		http.Redirect(w, r, "/page/"+pageGUID, http.StatusMovedPermanently)
	}
}
