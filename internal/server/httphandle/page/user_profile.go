package page

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"html/template"
	"net/http"

	"github.com/Anacardo89/lenic/internal/middleware"
	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/gorilla/mux"
)

type ProfilePage struct {
	Session *session.Session
	User    *models.User
	Posts   []*models.Post
	Follows string
}

func (h *PageHandler) UserProfile(w http.ResponseWriter, r *http.Request) {
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
	fDB, err := h.db.GetUserFollows(r.Context(), session.User.ID, uDB.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			fail("dberr: could not get follows", err, true, http.StatusInternalServerError, "internal error")
			return
		}
	}
	var pDB []*repo.Post
	if (session.User.ID == uDB.ID) ||
		(fDB != nil && fDB.FollowStatus == models.StatusAccepted.String()) {
		pDB, err = h.db.GetUserPosts(r.Context(), uDB.ID)
		if err != nil {
			fail("dberr: could not get posts", err, true, http.StatusInternalServerError, "internal error")
			return
		}
	} else {
		pDB, err = h.db.GetUserPublicPosts(r.Context(), uDB.ID)
		if err != nil {
			fail("dberr: could not get posts", err, true, http.StatusInternalServerError, "internal error")
			return
		}
	}
	// Response
	var (
		followStatus string
		posts        []*models.Post
	)
	u := models.FromDBUser(uDB)
	un := models.FromDBUserNotif(uDB)
	if fDB != nil {
		followStatus = fDB.FollowStatus
	} else {
		followStatus = models.StatusPending.String()
	}
	for _, post := range pDB {
		p := models.FromDBPost(post, *un)
		p.Content = template.HTML(p.RawContent)
		posts = append(posts, p)
	}
	pp := ProfilePage{
		Session: session,
		User:    u,
		Posts:   posts,
		Follows: followStatus,
	}
	t, err := template.ParseFiles("../frontend/templates/authorized/user-profile.html")
	if err != nil {
		fail("could not parse template", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	t.Execute(w, pp)
}
