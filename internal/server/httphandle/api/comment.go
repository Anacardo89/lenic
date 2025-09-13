package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Anacardo89/lenic/internal/helpers"
	"github.com/Anacardo89/lenic/internal/middleware"
	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/internal/server/wshandle"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AddCommentResponse struct {
	ID string `json:"id"`
}

// POST /action/post/{post_id}/comment
func (h *APIHandler) AddComment(w http.ResponseWriter, r *http.Request) {
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
	// Get Session
	session, ok := r.Context().Value(middleware.CtxKeySession).(*session.Session)
	if !ok {
		fail("session type mismatch", errors.New("session type mismatch"), true, http.StatusUnauthorized, "invalid session")
		return
	}
	// Input validation
	vars := mux.Vars(r)
	postID, err := uuid.Parse(vars["post_id"])
	if err != nil {
		fail("parsing post uuid from URL", err, true, http.StatusBadRequest, "invalid path")
		return
	}
	// DB operations
	c := repo.Comment{
		PostID:   postID,
		AuthorID: session.User.ID,
		Content:  r.FormValue("comment_text"),
	}
	if err := h.db.CreateComment(r.Context(), &c); err != nil {
		fail("dberr - could not insert comment", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Handle user mentions
	mentions := helpers.ParseAtString(c.Content)
	if len(mentions) > 0 {
		go func() {
			for _, mention := range mentions {
				mention = strings.TrimLeft(mention, "@")
				u, err := h.db.GetUserByUserName(h.ctx, mention)
				if err != nil {
					continue
				}
				ut := &repo.UserTag{
					UserID:      u.ID,
					TargetID:    c.ID,
					ResourceTpe: repo.ResourceComment.String(),
				}
				err = h.db.CreateUserTag(h.ctx, ut)
				if err != nil {
					fail("dberr - could not insert usertag", err, false, http.StatusInternalServerError, "")
					continue
				}
				wsMsg := wshandle.Message{
					FromUserName: session.User.UserName,
					Type:         "comment_tag",
					Msg:          " has tagged you in their comment",
					ResourceID:   c.ID.String(),
					ParentID:     c.PostID.String(),
				}
				h.wsHandler.HandleCommentTag(wsMsg, mention)
			}
		}()
	}
	// Response
	resp := AddCommentResponse{
		ID: c.ID.String(),
	}
	data, err := json.Marshal(&resp)
	if err != nil {
		fail("failed to marshal response body", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// PUT /action/post/{post_id}/comment/{comment_id}
func (h *APIHandler) EditComment(w http.ResponseWriter, r *http.Request) {
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
	cID, err := uuid.Parse(vars["comment_id"])
	if err != nil {
		fail("parsing comment uuid from URL", err, true, http.StatusBadRequest, "invalid path")
		return
	}
	if r.FormValue("comment") == "" {
		fail("comment is empty", errors.New("comment is empty"), true, http.StatusBadRequest, "comment cannot be empty")
		return
	}
	// DB operations
	c := repo.Comment{
		ID:      cID,
		Content: r.FormValue("comment"),
	}
	if err := h.db.UpdateComment(r.Context(), &c); err != nil {
		fail("dberr - could not update comment", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Handle user mentions
	mentions := helpers.ParseAtString(c.Content)
	if len(mentions) > 0 {
		go func() {
			for _, mention := range mentions {
				mention = strings.TrimLeft(mention, "@")
				u, err := h.db.GetUserByUserName(h.ctx, mention)
				if err != nil {
					continue
				}
				ut := &repo.UserTag{
					UserID:      u.ID,
					TargetID:    c.ID,
					ResourceTpe: repo.ResourceComment.String(),
				}
				err = h.db.CreateUserTag(h.ctx, ut)
				if err != nil {
					fail("dberr - could not insert usertag", err, false, http.StatusInternalServerError, "")
					continue
				}
				wsMsg := wshandle.Message{
					FromUserName: session.User.UserName,
					Type:         "comment_tag",
					Msg:          " has tagged you in their comment",
					ResourceID:   c.ID.String(),
					ParentID:     c.PostID.String(),
				}
				h.wsHandler.HandleCommentTag(wsMsg, mention)
			}
		}()
	}
	// Response
	w.WriteHeader(http.StatusOK)
}

// DELETE /action/post/{post_id}/comment/{comment_id}
func (h *APIHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
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
	// Input validation
	vars := mux.Vars(r)
	cID, err := uuid.Parse(vars["comment_id"])
	if err != nil {
		fail("parsing comment uuid from URL", err, true, http.StatusBadRequest, "invalid path")
		return
	}
	// DB operations
	c, err := h.db.DisableComment(r.Context(), cID)
	if err != nil {
		fail("dberr - could not disable comment", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// delete user mentions
	mentions := helpers.ParseAtString(c.Content)
	if len(mentions) > 0 {
		go func() {
			for _, mention := range mentions {
				mention = strings.TrimLeft(mention, "@")
				userDB, err := h.db.GetUserByUserName(h.ctx, mention)
				if err != nil {
					continue
				}
				err = h.db.DeleteUserTag(h.ctx, userDB.ID, cID)
				if err != nil {
					fail("dberr - could not delete usertag", err, false, http.StatusInternalServerError, "")
					continue
				}
			}
		}()
	}
	// Response
	w.WriteHeader(http.StatusOK)
}

// POST /action/post/{post_id}/comment/{comment_id}/up
func (h *APIHandler) RateCommentUp(w http.ResponseWriter, r *http.Request) {
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
	cID, err := uuid.Parse(vars["comment_id"])
	if err != nil {
		fail("parsing comment uuid from URL", err, true, http.StatusBadRequest, "invalid path")
		return
	}
	// DB operations
	err = h.db.RateCommentUp(r.Context(), cID, session.User.ID)
	if err != nil {
		fail("dberr - ould not update comment rating", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Response
	w.WriteHeader(http.StatusOK)
}

// POST /action/post/{post_id}/comment/{comment_id}/down
func (h *APIHandler) RateCommentDown(w http.ResponseWriter, r *http.Request) {
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
	cID, err := uuid.Parse(vars["comment_id"])
	if err != nil {
		fail("parsing comment uuid from URL", err, true, http.StatusBadRequest, "invalid path")
		return
	}
	// DB operations
	err = h.db.RateCommentDown(r.Context(), cID, session.User.ID)
	if err != nil {
		fail("dberr - ould not update comment rating", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Response
	w.WriteHeader(http.StatusOK)
}
