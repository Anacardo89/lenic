package actions

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

type LoginRequest struct {
	UserName     string `json:"user_name"`
	UserPassword string `json:"user_password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/login ", r.RemoteAddr)
	var (
		err      error
		loginReq LoginRequest
	)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error.Println("/action/login - Error reading body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &loginReq)
	if err != nil {
		logger.Error.Println("/action/login - Could not decode JSON: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	logger.Info.Printf("/action/login %s %s\n", r.RemoteAddr, loginReq.UserName)

	dbuser, err := orm.Da.GetUserByName(loginReq.UserName)
	if err == sql.ErrNoRows {
		http.Error(w, "User does not exist", http.StatusBadRequest)
		return
	}
	u := mapper.User(dbuser)
	if u.Active != 1 {
		if u.Active == 2 {
			http.Error(w, "User is blocked, contact the admin", http.StatusBadRequest)
			return
		}
		http.Error(w, "User is not active, check your mail", http.StatusBadRequest)
		return
	}
	u.Pass = loginReq.UserPassword
	if !auth.CheckPasswordHash(u.Pass, u.HashPass) {
		http.Error(w, "Password does not match", http.StatusBadRequest)
		return
	}
	usrSession := auth.CreateSession(w, r)
	auth.UpdateSession(usrSession.SessionId, u.Id)

	logger.Info.Println("OK - /action/login ", r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/logout ", r.RemoteAddr)
	session, err := auth.SessionStore.Get(r, "lenic")
	if err != nil {
		logger.Error.Println("/action/logout - Could not get session: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		logger.Error.Println("/action/logout - Could not save session: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info.Println("OK - /action/logout ", r.RemoteAddr)
	redirect.RedirIndex(w, r)
}
