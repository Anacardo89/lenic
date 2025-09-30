package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/Anacardo89/lenic/internal/helpers"
	"github.com/Anacardo89/lenic/internal/middleware"
	"github.com/Anacardo89/lenic/internal/repo"
	"github.com/Anacardo89/lenic/internal/server/wshandle"
	"github.com/Anacardo89/lenic/internal/session"
	"github.com/Anacardo89/lenic/pkg/fsops"
)

// POST /action/post
func (h *APIHandler) AddPost(w http.ResponseWriter, r *http.Request) {
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
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		fail("could not parse form", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	if r.FormValue("title") == "" || r.FormValue("content") == "" {
		fail("post without title or body", errors.New("invalid params"), true, http.StatusBadRequest, "post must contain a title and a body")
		return
	}
	// Make post
	isPublic := false
	if r.FormValue("is_public") == "1" {
		isPublic = true
	}
	pDB := &repo.Post{
		Title:    r.FormValue("title"),
		Content:  r.FormValue("content"),
		AuthorID: session.User.ID,
		IsPublic: isPublic,
	}
	// Handle image
	file, header, err := r.FormFile("post_image")
	if err != nil && err != http.ErrMissingFile {
		fail("could not get image", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	if file != nil && header != nil {
		fileExt := filepath.Ext(header.Filename)
		fileName := fsops.NameImg(16)
		fileName = fmt.Sprintf("%s.%s", fileName, fileExt)
		pDB.PostImage = fileName
		imgData, err := io.ReadAll(file)
		if err != nil {
			fail("could not read image data", err, true, http.StatusBadRequest, "invalid params")
		}
		fsops.SaveImg(imgData, fsops.PostImgPath, fileName)
	}
	// DB operations
	pID, err := h.db.CreatePost(h.ctx, pDB)
	if err != nil {
		fail("dberr: could not create post", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Handle user mentions
	mentions := helpers.ParseAtString(pDB.Content)
	titleMentions := helpers.ParseAtString(pDB.Title)
	mentions = append(mentions, titleMentions...)
	if len(mentions) > 0 {
		go func() {
			for _, mention := range mentions {
				mention = strings.TrimLeft(mention, "@")
				u, err := h.db.GetUserByUserName(h.ctx, mention)
				if err != nil {
					fail("dberr: could not get user", err, false, http.StatusInternalServerError, "")
					continue
				}
				ut := &repo.UserTag{
					UserID:      u.ID,
					TargetID:    pID,
					ResourceTpe: repo.ResourcePost.String(),
				}
				if err := h.db.CreateUserTag(h.ctx, ut); err != nil {
					fail("dberr: could not insert usertag", err, false, http.StatusInternalServerError, "")
					continue
				}
				wsMsg := wshandle.Message{
					FromUserName: session.User.Username,
					Type:         "post_tag",
					Msg:          " has tagged you in their post",
					ResourceID:   pID.String(),
				}
				h.wsHandler.HandlePostTag(wsMsg, mention)
			}
		}()
	}
	// Response
	w.WriteHeader(http.StatusCreated)
}

// PUT /action/post/{Post_GUID}
func (h *APIHandler) EditPost(w http.ResponseWriter, r *http.Request) {
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
	if err := r.ParseForm(); err != nil {
		fail("could not parse form", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	pID, err := uuid.Parse(vars["post_id"])
	if err != nil {
		fail("could not decode post_id", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	if r.FormValue("title") == "" || r.FormValue("content") == "" {
		fail("post without title or body", errors.New("invalid params"), true, http.StatusBadRequest, "post must contain a title and a body")
		return
	}
	// Make post
	isPublic := false
	if r.FormValue("is_public") == "1" {
		isPublic = true
	}
	p := repo.Post{
		ID:       pID,
		Title:    r.FormValue("title"),
		Content:  r.FormValue("content"),
		IsPublic: isPublic,
	}
	// DB operations
	if err := h.db.UpdatePost(h.ctx, &p); err != nil {
		fail("dberr: could not update post", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Handle user mentions
	mentions := helpers.ParseAtString(p.Content)
	titleMentions := helpers.ParseAtString(p.Title)
	mentions = append(mentions, titleMentions...)
	if len(mentions) > 0 {
		go func() {
			for _, mention := range mentions {
				mention = strings.TrimLeft(mention, "@")
				u, err := h.db.GetUserByUserName(h.ctx, mention)
				if err != nil {
					fail("dberr: could not get user", err, false, http.StatusInternalServerError, "")
					continue
				}
				ut := &repo.UserTag{
					UserID:      u.ID,
					TargetID:    pID,
					ResourceTpe: repo.ResourcePost.String(),
				}
				if err := h.db.CreateUserTag(h.ctx, ut); err != nil {
					fail("dberr: could not insert usertag", err, false, http.StatusInternalServerError, "")
					continue
				}
				wsMsg := wshandle.Message{
					FromUserName: session.User.Username,
					Type:         "post_tag",
					Msg:          " has tagged you in their post",
					ResourceID:   pID.String(),
				}
				h.wsHandler.HandlePostTag(wsMsg, mention)
			}
		}()
	}
	// Response
	w.WriteHeader(http.StatusOK)
}

// DELETE /action/post/{Post_GUID}
func (h *APIHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
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
	pID, err := uuid.Parse(vars["post_id"])
	if err != nil {
		fail("could not decode post_id", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// DB operations
	p, err := h.db.DisablePost(h.ctx, pID)
	if err != nil {
		fail("dberr: could not disable post", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// Delete user mentions
	mentions := helpers.ParseAtString(p.Content)
	titleMentions := helpers.ParseAtString(p.Title)
	mentions = append(mentions, titleMentions...)
	if len(mentions) > 0 {
		go func() {
			for _, mention := range mentions {
				mention = strings.TrimLeft(mention, "@")
				u, err := h.db.GetUserByUserName(h.ctx, mention)
				if err != nil {
					fail("dberr: could not get user", err, false, http.StatusInternalServerError, "")
					continue
				}
				if err := h.db.DeleteUserTag(h.ctx, u.ID, pID); err != nil {
					fail("dberr: could not delete usertag", err, false, http.StatusInternalServerError, "")
					continue
				}
			}
		}()
	}
	// Response
	w.WriteHeader(http.StatusOK)
}

// POST /action/post/{Post_GUID}/up
func (h *APIHandler) RatePostUp(w http.ResponseWriter, r *http.Request) {
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
	if err := h.db.RatePostUp(h.ctx, pID, session.User.ID); err != nil {
		fail("dberr: could not update post rating", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// Response
	w.WriteHeader(http.StatusOK)
}

// POST /action/post/{Post_GUID}/down
func (h *APIHandler) RatePostDown(w http.ResponseWriter, r *http.Request) {
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
	err = h.db.RatePostDown(h.ctx, pID, session.User.ID)
	if err != nil {
		fail("dberr: could not update post rating", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// Response
	w.WriteHeader(http.StatusOK)
}
