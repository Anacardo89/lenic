package page

import (
	"encoding/base64"
	"errors"
	"html/template"
	"net/http"

	"github.com/Anacardo89/lenic/internal/middleware"
	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/gorilla/mux"
)

type FollowersPage struct {
	Session   *session.Session
	User      *models.User
	Followers []*models.User
}

type FollowingPage struct {
	Session   *session.Session
	User      *models.User
	Following []*models.User
}

func (h *PageHandler) UserFollowers(w http.ResponseWriter, r *http.Request) {
	// Error Handling
	fail := func(logMsg string, e error, writeError bool, status int, outMsg string) {
		h.log.Error(logMsg, "error", e,
			"status_code", status,
			"method", r.Method,
			"path", r.URL.Path,
			"client_ip", r.RemoteAddr,
		)
		if writeError {
			http.Error(w, outMsg, status)
		}
	}
	//

	// Execution
	// Get session
	session, ok := r.Context().Value(middleware.CtxKeySession).(*session.Session)
	if !ok {
		fail("session type mismatch", errors.New("session type mismatch"), true, http.StatusUnauthorized, "invalid session")
		return
	}
	// Input validation
	vars := mux.Vars(r)
	bytes, err := base64.URLEncoding.DecodeString(vars["encoded_username"])
	if err != nil {
		fail("could not decode user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	username := string(bytes)
	// DB operations
	uDB, err := h.db.GetUserByUserName(r.Context(), username)
	if err != nil {
		fail("dberr: could not get user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	fsDB, err := h.db.GetFollowers(r.Context(), uDB.ID)
	if err != nil {
		fail("dberr: could not get followers", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Response
	u := models.FromDBUser(uDB)
	fp := FollowersPage{
		Session: session,
		User:    u,
	}
	for _, f := range fsDB {
		uDB, err := h.db.GetUserByID(r.Context(), f.FollowerID)
		if err != nil {
			fail("dberr: could not get user", err, true, http.StatusInternalServerError, "internal error")
			return
		}
		u := models.FromDBUser(uDB)
		fp.Followers = append(fp.Followers, u)
	}
	t, err := template.ParseFiles("../frontend/templates/authorized/followers.html")
	if err != nil {
		fail("could not parse template", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	t.Execute(w, fp)
}

func (h *PageHandler) UserFollowing(w http.ResponseWriter, r *http.Request) {
	// Error Handling
	fail := func(logMsg string, e error, writeError bool, status int, outMsg string) {
		h.log.Error(logMsg, "error", e,
			"status_code", status,
			"method", r.Method,
			"path", r.URL.Path,
			"client_ip", r.RemoteAddr,
		)
		if writeError {
			http.Error(w, outMsg, status)
		}
	}
	//

	// Execution
	// Get session
	session, ok := r.Context().Value(middleware.CtxKeySession).(*session.Session)
	if !ok {
		fail("session type mismatch", errors.New("session type mismatch"), true, http.StatusUnauthorized, "invalid session")
		return
	}
	// Input validation
	vars := mux.Vars(r)
	bytes, err := base64.URLEncoding.DecodeString(vars["encoded_username"])
	if err != nil {
		fail("could not decode user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	username := string(bytes)
	// DB operations
	uDB, err := h.db.GetUserByUserName(h.ctx, username)
	if err != nil {
		fail("dberr: could not get user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	fsDB, err := h.db.GetFollowing(h.ctx, uDB.ID)
	if err != nil {
		fail("dberr: could not get following", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Response
	u := models.FromDBUser(uDB)
	fp := FollowingPage{
		Session: session,
		User:    u,
	}
	for _, f := range fsDB {
		dbUser, err := h.db.GetUserByID(h.ctx, f.FollowedID)
		if err != nil {
			fail("dberr: could not get user", err, true, http.StatusInternalServerError, "internal error")
			return
		}
		u := models.FromDBUser(dbUser)
		fp.Following = append(fp.Following, u)
	}
	t, err := template.ParseFiles("../frontend/templates/authorized/following.html")
	if err != nil {
		fail("could not parse template", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	t.Execute(w, fp)
}
