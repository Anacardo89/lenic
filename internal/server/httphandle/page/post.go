package page

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/Anacardo89/lenic/internal/middleware"
	"github.com/Anacardo89/lenic/internal/models"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PostPage struct {
	Session *session.Session
	Post    *models.Post
}

func (h *PageHandler) NewPost(w http.ResponseWriter, r *http.Request) {
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
	// Response
	postp := PostPage{
		Session: session,
		Post:    &models.Post{},
	}
	t, err := template.ParseFiles("templates/authorized/newPost.html")
	if err != nil {
		fail("could not parse template", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	t.Execute(w, postp)
}

// /post/{post_id}
func (h *PageHandler) Post(w http.ResponseWriter, r *http.Request) {
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
	pID, err := uuid.Parse(vars["post_id"])
	if err != nil {
		fail("could not decode post_id", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// DB operations
	pDB, err := h.db.GetPostForPage(r.Context(), pID, session.User.ID)
	if err != nil {
		fail("dberr: could not get post", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// Response
	u := models.FromDBUserNotif(&pDB.Author)
	p := models.FromDBPost(&pDB.Post, *u)
	p.UserRating = pDB.UserRating
	p.Content = template.HTML(p.RawContent)
	var comments []*models.Comment
	for _, comment := range pDB.Comments {
		cu := models.FromDBUserNotif(&comment.Author)
		c := models.FromDBComment(&comment.Comment, *cu)
		c.UserRating = comment.UserRating
		comments = append(comments, c)
	}
	p.Comments = comments
	pp := PostPage{
		Post:    p,
		Session: session,
	}
	t, err := template.ParseFiles("../frontend/templates/authorized/post.html")
	if err != nil {
		fail("could not parse template", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	t.Execute(w, pp)
}
