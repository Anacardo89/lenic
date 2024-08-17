package actions

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Anacardo89/tpsi25_blog/auth"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/internal/rabbit"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/Anacardo89/tpsi25_blog/pkg/rabbitmq"
)

type RegisterData struct {
	Email string `json:"email"`
	User  string `json:"user"`
	Link  string `json:"link"`
}

func isValidInput(input string) bool {
	return !strings.Contains(input, ";")
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Parse Form
	err := r.ParseForm()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	var u = &presentation.User{
		UserName:  r.FormValue("user_name"),
		UserEmail: r.FormValue("user_email"),
		UserPass:  r.FormValue("user_password"),
		Active:    0,
	}
	pass2 := r.FormValue("user_password2")
	if u.UserPass != pass2 {
		RedirectToError(w, r, "Password strings don't match")
		return
	}
	if !isValidInput(u.UserName) || !isValidInput(u.UserEmail) || !isValidInput(u.UserPass) {
		RedirectToError(w, r, "Invalid character in form")
		return
	}

	// Check if UserName or Email in use
	_, err = orm.Da.GetUserByName(u.UserName)
	if err != sql.ErrNoRows {
		RedirectToError(w, r, "User already exists")
		return
	}
	_, err = orm.Da.GetUserByEmail(u.UserEmail)
	if err != sql.ErrNoRows {
		RedirectToError(w, r, "Email already exists")
		return
	}

	// Password Hashing
	u.HashedPass, err = auth.HashPassword(u.UserPass)
	if err != nil {
		logger.Error.Println(err)
		return
	}

	// Send Regsiter Mail to Queue
	regData := RegisterData{
		Email: u.UserEmail,
		User:  u.UserName,
		Link:  generateActiveLink(u.UserName),
	}
	data, err := json.Marshal(regData)
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}

	err = rabbit.MQSendRegisterMail(rabbitmq.RMQ, rabbitmq.RCh, data)
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}

	// Insert User in DB
	dbuser := mapper.UserToDB(u)
	err = orm.Da.CreateUser(dbuser)
	if err != nil {
		logger.Error.Println(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func generateActiveLink(user string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(user))
	return "https://192.168.200.205:8082/activate/" + encoded
}