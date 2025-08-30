package api

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/gorilla/mux"
)

// GET /action/search/user
func (h *APIHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	username := queryParams.Get("username")

	dbusers, err := h.db.SearchUsersByUserName(h.ctx, username)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error.Printf("GET /action/search/user %s - Could not get users: %s\n", username, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if dbusers == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	var users []models.UserNotif
	for _, dbuser := range dbusers {
		u := models.FromDBUserNotif(&dbuser)
		users = append(users, *u)
	}

	data, err := json.Marshal(users)
	if err != nil {
		logger.Error.Printf("GET /action/search/user %s - Could not marshal users: %s\n", username, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// POST /action/user/{user_encoded}/follow
func (h *APIHandler) RequestFollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/follow - Could not decode user: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	username := string(bytes)

	dbuser, err := h.db.GetUserByUserName(h.ctx, username)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/follow - Could not get user: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session := h.sessionStore.ValidateSession(w, r)

	err = h.db.FollowUser(h.ctx, session.User.ID, dbuser.ID)
	if err != nil {
		logger.Error.Printf("POST /action/user/%s/follow - Could not follow: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DELETE /action/user/{user_encoded}/unfollow
func (h *APIHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("DELETE /action/user/%s/unfollow - Could not decode user: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	username := string(bytes)

	dbuser, err := h.db.GetUserByUserName(h.ctx, username)
	if err != nil {
		logger.Error.Printf("DELETE /action/user/%s/unfollow - Could not get user: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	requesterName := queryParams.Get("requester")

	dbrequester, err := h.db.GetUserByUserName(h.ctx, requesterName)
	if err != nil {
		logger.Error.Printf("DELETE /action/user/%s/unfollow - Could not get dbrequester: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.UnfollowUser(h.ctx, dbrequester.ID, dbuser.ID)
	if err != nil {
		logger.Error.Printf("DELETE /action/user/%s/unfollow - Could not unfollow: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbnotif, err := h.db.GetFollowNotification(h.ctx, dbuser.ID, dbrequester.ID)
	if err != nil {
		logger.Error.Printf("DELETE /action/user/%s/unfollow - Could not get notif: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.DeleteNotification(h.ctx, dbnotif.ID)
	if err != nil {
		logger.Error.Printf("DELETE /action/user/%s/unfollow - Could not delete notif: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// PUT /action/user/{user_encoded}/accept
func (h *APIHandler) AcceptFollowRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not decode user: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	username := string(bytes)

	err = r.ParseForm()
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not parse form: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	requesterName := r.FormValue("requester")

	dbuser, err := h.db.GetUserByUserName(h.ctx, username)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not decode user: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbrequester, err := h.db.GetUserByUserName(h.ctx, requesterName)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not decode user: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.AcceptFollow(h.ctx, dbrequester.ID, dbuser.ID)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not accept follow: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbnotif, err := h.db.GetFollowNotification(h.ctx, dbuser.ID, dbrequester.ID)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not get notif: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.DeleteNotification(h.ctx, dbnotif.ID)
	if err != nil {
		logger.Error.Printf("PUT /action/user/%s/accept - Could not delete notif: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
