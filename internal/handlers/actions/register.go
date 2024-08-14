package actions

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Anacardo89/tpsi25_blog/auth"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/internal/rabbit"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

type RegisterData struct {
	Email string `json:"email"`
	User  string `json:"user"`
	Link  string `json:"link"`
}

func isValidInput(input string) bool {
	return !strings.Contains(input, ";")
}

func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	// Parse Form
	err := r.ParseForm()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	var userReg = &presentation.User{
		UserName:  r.FormValue("user_name"),
		UserEmail: r.FormValue("user_email"),
		UserPass:  r.FormValue("user_password"),
	}
	pass2 := r.FormValue("user_password2")
	if userReg.UserPass != pass2 {
		cookie := http.Cookie{Name: "errormsg",
			Value:    "Password strings don't match",
			Expires:  time.Now().Add(60 * time.Second),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/error", http.StatusMovedPermanently)
		return
	}
	if !isValidInput(userReg.UserName) || !isValidInput(userReg.UserEmail) || !isValidInput(userReg.UserPass) {
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

	// Check if UserName or Email in use
	dbUser := database.User{}
	err = db.Dbase.QueryRow(query.SelectUserByName,
		userReg.UserName).
		Scan(dbUser.UserName)
	if err != sql.ErrNoRows {
		cookie := http.Cookie{Name: "errormsg",
			Value:    "User already exists",
			Expires:  time.Now().Add(60 * time.Second),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/error", http.StatusMovedPermanently)
		return
	}
	err = db.Dbase.QueryRow(query.SelectUserByEmail,
		userReg.UserEmail).
		Scan(dbUser.UserEmail)
	if err != sql.ErrNoRows {
		cookie := http.Cookie{Name: "errormsg",
			Value:    "Email already exists",
			Expires:  time.Now().Add(60 * time.Second),
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/error", http.StatusMovedPermanently)
		return
	}

	// Password Hashing
	userReg.HashedPass, err = auth.HashPassword(userReg.UserPass)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	// Send Regsiter Mail to Queue
	regData := RegisterData{
		Email: userReg.UserEmail,
		User:  userReg.UserName,
		Link:  generateActiveLink(userReg.UserName),
	}

	var mbuf bytes.Buffer
	regData.Email = mbuf.String()
	mbuf.Reset()
	regData.User = mbuf.String()
	mbuf.Reset()
	regData.Link = mbuf.String()
	data, err := json.Marshal(regData)
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}

	err = rabbit.RabbitMQ.MQSendRegisterMail(data)
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}

	// Insert User in DB
	_, err = db.Dbase.Exec(query.InsertUser,
		userReg.UserName, userReg.UserEmail, userReg.HashedPass, 0)
	if err != nil {
		logger.Error.Println(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func generateActiveLink(user string) string {
	return "https://192.168.200.205:8082/activate/" + user
}
