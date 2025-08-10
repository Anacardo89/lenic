package server

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/internal/helpers"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mqmodel"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/internal/rabbit"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/Anacardo89/tpsi25_blog/pkg/rabbitmq"
	"github.com/gorilla/mux"
)

// /action/register
func (s *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/register ", r.RemoteAddr)
	// Parse Form
	err := r.ParseForm()
	if err != nil {
		logger.Error.Println("/action/register - Could not parse Form: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	var u = &presentation.User{
		UserName:   r.FormValue("user_name"),
		Email:      r.FormValue("user_email"),
		Pass:       r.FormValue("user_password"),
		ProfilePic: "",
		Active:     0,
	}
	pass2 := r.FormValue("user_password2")
	if u.Pass != pass2 {
		redirect.RedirectToError(w, r, "Password strings don't match")
		return
	}

	// Check if UserName or Email in use
	_, err = orm.Da.GetUserByName(u.UserName)
	if err != sql.ErrNoRows {
		logger.Error.Println("/action/register - Could not get user by name: ", err)
		redirect.RedirectToError(w, r, "User already exists")
		return
	}
	_, err = orm.Da.GetUserByEmail(u.Email)
	if err != sql.ErrNoRows {
		logger.Error.Println("/action/register - Could not get user by mail: ", err)
		redirect.RedirectToError(w, r, "Email already exists")
		return
	}

	// Password Hashing
	u.HashPass, err = auth.HashPassword(u.Pass)
	if err != nil {
		logger.Error.Println("/action/register - Could not hash password: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	// Send Regsiter Mail to Queue
	msg := mqmodel.Register{
		Email: u.Email,
		User:  u.UserName,
		Link:  helpers.MakeActivateUserLink(u.UserName),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error.Println("/action/register - Could not marshal JSON: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	err = rabbit.MQSendRegisterMail(rabbitmq.RMQ, rabbitmq.RCh, data)
	if err != nil {
		logger.Error.Println("/action/register - Could not send MQ msg: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	// Insert User in DB
	dbuser := mapper.UserToDB(u)
	dbuser.ProfilePicExt = ""
	err = orm.Da.CreateUser(dbuser)
	if err != nil {
		logger.Error.Println("/action/register - Could not create user: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	logger.Info.Printf("OK - /action/register %s %s\n", r.RemoteAddr, dbuser.UserName)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

// /action/activate
func (s *Server) ActivateUser(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/activate ", r.RemoteAddr)
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Println("/action/activate - Could not decode user: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("/action/activate %s %s\n", r.RemoteAddr, userName)
	err = orm.Da.SetUserAsActive(userName)
	if err != nil {
		logger.Error.Println("/action/activate - Could not activate user: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	logger.Info.Printf("OK - /action/activate %s %s\n", r.RemoteAddr, userName)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
