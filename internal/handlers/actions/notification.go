package actions

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

func GetNotifs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("GET /action/user/%s/notifications %s\n", encoded, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("GET /action/user/%s/notifications %s %s\n", encoded, r.RemoteAddr, userName)

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not get user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	queryParams := r.URL.Query()
	offset := queryParams.Get("offset")
	offsetint, err := strconv.Atoi(offset)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not parse offset to int: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	limit := queryParams.Get("limit")
	limitint, err := strconv.Atoi(limit)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not parse limit to int: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	dbnotifs, err := orm.Da.GetNotificationsByUser(dbuser.Id, limitint, offsetint)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not get notifs: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	var notifs []*presentation.Notification
	for _, dbnotif := range dbnotifs {
		dbfromuser, err := orm.Da.GetUserByID(dbnotif.FromUserId)
		if err != nil {
			logger.Error.Printf("GET /action/user/%s/notifications - Could not get notifs: %s\n", encoded, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
		n := mapper.Notification(dbnotif, dbuser.UserName, dbfromuser.UserName)
		notifs = append(notifs, n)
	}

	data, err := json.Marshal(notifs)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not marshal notifs: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
