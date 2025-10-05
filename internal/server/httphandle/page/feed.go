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

type FeedPage struct {
	Session *session.Session
	Posts   []*models.Post
}

func (h *PageHandler) Feed(w http.ResponseWriter, r *http.Request) {
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
	postsDB, err := h.db.GetFeed(r.Context(), username)
	if err != nil {
		fail("dberr: could not get feed", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// Response
	feed := FeedPage{
		Session: session,
	}
	for _, p := range postsDB {
		uDB, err := h.db.GetUserByID(r.Context(), p.AuthorID)
		if err != nil {
			fail("dberr: could not get user", err, true, http.StatusInternalServerError, "internal error")
			return
		}
		u := models.FromDBUserNotif(uDB)
		post := models.FromDBPost(p, *u)
		post.Content = template.HTML(post.RawContent)
		feed.Posts = append(feed.Posts, post)
	}
	t, err := template.ParseFiles("../frontend/templates/authorized/feed.html")
	if err != nil {
		fail("could not parse template", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	t.Execute(w, feed)
}
