package api

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Anacardo89/lenic/internal/helpers"
	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/models/mqmodel"
	"github.com/Anacardo89/lenic/internal/rabbit"
	"github.com/Anacardo89/lenic/internal/server/httphandle/redirect"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/Anacardo89/lenic/pkg/rabbitmq"
	"github.com/gorilla/mux"
)

// /action/register
func (h *APIHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/register ", r.RemoteAddr)
	// Parse Form
	err := r.ParseForm()
	if err != nil {
		logger.Error.Println("/action/register - Could not parse Form: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	var u = &models.User{
		UserName:   r.FormValue("username"),
		Email:      r.FormValue("email"),
		Pass:       r.FormValue("password"),
		ProfilePic: "",
	}
	pass2 := r.FormValue("user_password2")
	if u.Pass != pass2 {
		redirect.RedirectToError(w, r, "Password strings don't match")
		return
	}

	// Check if UserName or Email in use
	_, err = h.db.GetUserByUserName(h.ctx, u.UserName)
	if err != sql.ErrNoRows {
		logger.Error.Println("/action/register - Could not get user by name: ", err)
		redirect.RedirectToError(w, r, "Username already exists")
		return
	}
	_, err = h.db.GetUserByEmail(h.ctx, u.Email)
	if err != sql.ErrNoRows {
		logger.Error.Println("/action/register - Could not get user by mail: ", err)
		redirect.RedirectToError(w, r, "Email already exists")
		return
	}

	// Password Hashing
	u.PasswordHash, err = helpers.HashPassword(u.Pass)
	if err != nil {
		logger.Error.Println("/action/register - Could not hash password: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	// Send Regsiter Mail to Queue
	msg := mqmodel.Register{
		Email: u.Email,
		User:  u.UserName,
		Link:  helpers.MakeActivateUserLink(h.cfg.Host, h.cfg.HTTPSPort, u.UserName),
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
	uDB := models.ToDBUser(u)
	_, err = h.db.CreateUser(h.ctx, uDB)
	if err != nil {
		logger.Error.Println("/action/register - Could not create user: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	logger.Info.Printf("OK - /action/register %s %s\n", r.RemoteAddr, uDB.UserName)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

// /action/activate
func (h *APIHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
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
	err = h.db.SetUserActive(h.ctx, userName)
	if err != nil {
		logger.Error.Println("/action/activate - Could not activate user: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	logger.Info.Printf("OK - /action/activate %s %s\n", r.RemoteAddr, userName)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}
