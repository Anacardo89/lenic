package actions

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

// GET /action/user/{user_encoded}/notifications
func GetNotifs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("GET /action/user/%s/notifications %s\n", encoded, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not decode user: %s\n", encoded, err)
		return
	}
	userName := string(bytes)
	logger.Info.Printf("GET /action/user/%s/notifications %s %s\n", encoded, r.RemoteAddr, userName)

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not get user: %s\n", encoded, err)
		return
	}
	u := mapper.UserNotif(dbuser)

	queryParams := r.URL.Query()
	offset := queryParams.Get("offset")
	offsetint, err := strconv.Atoi(offset)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not parse offset to int: %s\n", encoded, err)
		return
	}

	limit := queryParams.Get("limit")
	limitint, err := strconv.Atoi(limit)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not parse limit to int: %s\n", encoded, err)
		return
	}

	dbnotifs, err := orm.Da.GetNotificationsByUser(dbuser.Id, limitint, offsetint)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not get notifs: %s\n", encoded, err)
		return
	}

	var notifs []*presentation.Notification
	for _, dbnotif := range dbnotifs {
		dbfromuser, err := orm.Da.GetUserByID(dbnotif.FromUserId)
		if err != nil {
			logger.Error.Printf("GET /action/user/%s/notifications - Could not get user: %s\n", encoded, err)
			return
		}
		from_u := mapper.UserNotif(dbfromuser)
		n := mapper.Notification(dbnotif, *u, *from_u)
		notifs = append(notifs, n)
	}

	data, err := json.Marshal(notifs)
	if err != nil {
		logger.Error.Printf("GET /action/user/%s/notifications - Could not marshal notifs: %s\n", encoded, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// PUT /action/user/{user_encoded}/notifications/{notif_id}/read
func UpdateNotif(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	notif_id := vars["notif_id"]
	logger.Info.Printf("PUT /action/user/%s/notifications/%s/read %s\n", encoded, notif_id, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/notifications/%s/read - Could not decode user: %s\n", encoded, notif_id, err)
		return
	}
	userName := string(bytes)
	logger.Info.Printf("PUT /action/user/%s/notifications/%s/read %s %s\n", encoded, notif_id, r.RemoteAddr, userName)

	notif_id_int, err := strconv.Atoi(notif_id)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/notifications/%s/read - Could not parse notif_id to int: %s\n", encoded, notif_id, err)
		return
	}

	err = orm.Da.UpdateNotificationRead(notif_id_int)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/notifications/%s/read - Could not update notif: %s\n", encoded, notif_id, err)
		return
	}
}
