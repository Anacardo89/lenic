package pages

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile("../templates/login.html")
	if err != nil {
		logger.Error.Println(err)
	}
	fmt.Fprint(w, string(body))
}

func Register(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile("../templates/register.html")
	if err != nil {
		logger.Error.Println(err)
	}
	fmt.Fprint(w, string(body))
}

func ActivateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Println(err)
	}
	userName := string(bytes)
	err = orm.Da.SetUserAsActive(userName)
	if err != nil {
		logger.Error.Println(err)
	}
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	body, err := os.ReadFile("../templates/forgot-password.html")
	if err != nil {
		logger.Error.Println(err)
	}
	fmt.Fprint(w, string(body))
}
