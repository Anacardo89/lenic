package actions

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mqmodel"
	"github.com/Anacardo89/tpsi25_blog/internal/rabbit"
	"github.com/Anacardo89/tpsi25_blog/internal/server"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/Anacardo89/tpsi25_blog/pkg/rabbitmq"
)

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/forgot-password ", r.RemoteAddr)
	// Parse Form
	err := r.ParseForm()
	if err != nil {
		logger.Error.Println("/action/forgot-password - Could not parse Form: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mail := r.FormValue("user_email")
	logger.Info.Printf("/action/forgot-password %s %s\n", r.RemoteAddr, mail)
	// Get user from DB
	dbuser, err := orm.Da.GetUserByEmail(mail)
	if err == sql.ErrNoRows {
		http.Error(w, "No user with that email", http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateToken(64)
	if err != nil {
		logger.Error.Println("/action/forgot-password - Could not generate token: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t := &database.Token{
		Token:  token,
		UserId: dbuser.Id,
	}

	err = orm.Da.CreateToken(t)
	if err != nil {
		logger.Error.Println("/action/forgot-password - Could not create db token: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg := mqmodel.PasswordRecover{
		Email: dbuser.Email,
		User:  dbuser.UserName,
		Link:  makePasswordRecoverMail(dbuser.UserName, token),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error.Println("/action/forgot-password - Could not marshal JSON: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rabbit.MQSendPasswordRecoveryMail(rabbitmq.RMQ, rabbitmq.RCh, data)
	if err != nil {
		logger.Error.Println("/action/forgot-password - Could not send MQ msg: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info.Println("OK - /action/forgot-password ", r.RemoteAddr)
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func makePasswordRecoverMail(user string, token string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(user))
	return fmt.Sprintf("https://%s:%s/recover-password/%s?token=%s", server.Server.Host, server.Server.HttpsPORT, encoded, token)
}
