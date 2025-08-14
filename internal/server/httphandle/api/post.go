package api

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/helpers"
	"github.com/Anacardo89/lenic/internal/server/wshandle"
	"github.com/Anacardo89/lenic/pkg/fsops"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/google/uuid"

	"github.com/gorilla/mux"
)

// POST /action/post
func (h *APIHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/post ", r.RemoteAddr)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error.Println("/action/post - Could not parse Form: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	session := h.sessionStore.ValidateSession(w, r)

	isPublic := false
	isPublicStr := r.FormValue("is_public")
	if isPublicStr == "1" {
		isPublic = true
	}

	pDB := db.Post{
		Title:    r.FormValue("title"),
		Content:  r.FormValue("content"),
		AuthorID: session.User.ID,
		IsPublic: isPublic,
	}
	if pDB.Title == "" || pDB.Content == "" {
		http.Error(w, "Post must contain a title and a body", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("post_image")
	if err != nil {
		if err != http.ErrMissingFile {
			logger.Error.Println("/action/post - Could not get image: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pID, err := h.db.CreatePost(h.ctx, &pDB)
		if err != nil {
			logger.Error.Println("/action/post - Could not create post: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pDB, err = h.db.GetPost(h.ctx, pID)
		if err != nil {
			logger.Error.Println("/action/post - Could not get post: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mentions := helpers.ParseAtString(pDB.Content)
		titleMentions := helpers.ParseAtString(pDB.Title)
		mentions = append(mentions, titleMentions...)
		if len(mentions) > 0 {
			for _, mention := range mentions {
				mention = strings.TrimLeft(mention, "@")
				userDB, err := h.db.GetUserByUserName(h.ctx, mention)
				ut := &db.UserTag{
					UserID:      userDB.ID,
					TargetID:    pDB.ID,
					ResourceTpe: db.ResourcePost.String(),
				}
				err = h.db.CreateUserTag(h.ctx, ut)
				if err != nil {
					logger.Error.Printf("POST /action/post/%s/comment - Could not create UserTag: %s\n", postID, err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				wsMsg := wshandle.Message{
					FromUserName: session.User.UserName,
					Type:         "post_tag",
					Msg:          " has tagged you in their post",
					ResourceID:   pDB.ID.String(),
				}
				h.wsHandler.HandleCommentTag(wsMsg, mention)
			}
		}

		w.WriteHeader(http.StatusCreated)
		return
	}

	// Handle uploaded image
	fileExt := filepath.Ext(header.Filename)
	fileName := fsops.NameImg(16)
	fileName = fmt.Sprintf("%s.%s", fileName, fileExt)
	pDB.PostImage = fileName
	imgData, err := io.ReadAll(file)
	if err != nil {
		logger.Error.Println("/action/post - Could not read image data: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fsops.SaveImg(imgData, fsops.PostImgPath, fileName)

	// Insert post with image data
	pID, err := h.db.CreatePost(h.ctx, &pDB)
	if err != nil {
		logger.Error.Println("/action/post - Could not not create post: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pDB, err = h.db.GetPost(h.ctx, pID)
	if err != nil {
		logger.Error.Println("/action/post - Could not get post: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mentions := helpers.ParseAtString(pDB.Content)
	titleMentions := helpers.ParseAtString(pDB.Title)
	mentions = append(mentions, titleMentions...)
	if len(mentions) > 0 {
		for _, mention := range mentions {
			mention = strings.TrimLeft(mention, "@")
			userDB, err := h.db.GetUserByUserName(h.ctx, mention)
			ut := &db.UserTag{
				UserID:      userDB.ID,
				TargetID:    pDB.ID,
				ResourceTpe: db.ResourcePost.String(),
			}
			err = h.db.CreateUserTag(h.ctx, ut)
			if err != nil {
				logger.Error.Printf("POST /action/post/%s/comment - Could not create UserTag: %s\n", postID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			wsMsg := wshandle.Message{
				FromUserName: session.User.UserName,
				Type:         "post_tag",
				Msg:          " has tagged you in their post",
				ResourceID:   pDB.ID.String(),
			}
			h.wsHandler.HandleCommentTag(wsMsg, mention)
		}
	}

	logger.Info.Println("OK - /action/post ", r.RemoteAddr)
	w.WriteHeader(http.StatusCreated)
}

// PUT /action/post/{Post_GUID}
func (h *APIHandler) EditPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pIDstr := vars["post_id"]
	logger.Info.Printf("PUT /action/post/%s %s\n", pIDstr, r.RemoteAddr)
	session := h.sessionStore.ValidateSession(w, r)

	err := r.ParseForm()
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s - Could not parse form: %s\n", pIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pID, err := uuid.Parse(pIDstr)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s - Could not convert id to string: %s\n", pIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isPublic := false
	isPublicStr := r.FormValue("is_public")
	if isPublicStr == "1" {
		isPublic = true
	}

	p := db.Post{
		ID:       pID,
		Title:    r.FormValue("title"),
		Content:  r.FormValue("content"),
		IsPublic: isPublic,
	}
	if p.Content == "" || p.Title == "" {
		http.Error(w, "All form fields must be filled out", http.StatusBadRequest)
		return
	}

	err = h.db.UpdatePost(h.ctx, &p)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s - Could not update post: %s\n", pIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mentions := helpers.ParseAtString(p.Content)
	titleMentions := helpers.ParseAtString(pDB.Title)
	mentions = append(mentions, titleMentions...)
	if len(mentions) > 0 {
		for _, mention := range mentions {
			mention = strings.TrimLeft(mention, "@")
			userDB, err := h.db.GetUserByUserName(h.ctx, mention)
			ut := &db.UserTag{
				UserID:      userDB.ID,
				TargetID:    p.ID,
				ResourceTpe: db.ResourcePost.String(),
			}
			err = h.db.CreateUserTag(h.ctx, ut)
			if err != nil {
				logger.Error.Printf("POST /action/post/%s/comment - Could not create UserTag: %s\n", pIDstr, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			wsMsg := wshandle.Message{
				FromUserName: session.User.UserName,
				Type:         "post_tag",
				Msg:          " has tagged you in their post",
				ResourceID:   p.ID.String(),
			}
			h.wsHandler.HandleCommentTag(wsMsg, mention)
		}
	}

	logger.Info.Printf("OK - PUT /action/post/%s %s\n", pIDstr, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

// DELETE /action/post/{Post_GUID}
func (h *APIHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pIDstr := vars["post_id"]
	logger.Info.Printf("DELETE /action/post/%s %s\n", pIDstr, r.RemoteAddr)

	pID, err := uuid.Parse(pIDstr)
	if err != nil {
		logger.Error.Printf("DELETE /action/post/%s - Could not convert id to string: %s\n", pIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pDB, err := h.db.GetPost(h.ctx, pID)
	if err != nil {
		logger.Error.Printf("DELETE /action/post/%s - Could not get comment: %s\n", pIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.DisablePost(h.ctx, pID)
	if err != nil {
		logger.Error.Printf("DELETE /action/post/%s - Could not update comment: %s\n", pIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mentions := helpers.ParseAtString(pDB.Content)
	titleMentions := helpers.ParseAtString(pDB.Title)
	mentions = append(mentions, titleMentions...)
	if len(mentions) > 0 {
		for _, mention := range mentions {
			mention = strings.TrimLeft(mention, "@")
			uDB, err := h.db.GetUserByUserName(h.ctx, mention)
			if err != nil {
				logger.Error.Printf("DELETE /action/post/%s - Could not get tag By Id: %s\n", pIDstr, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = h.db.DeleteUserTag(h.ctx, uDB.ID, pDB.ID)
			if err != nil {
				logger.Error.Printf("DELETE /action/post/%s - Could not delete tag By Id: %s\n", pIDstr, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	logger.Info.Printf("OK - DELETE /action/post/%s %s\n", pIDstr, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

// POST /action/post/{Post_GUID}/up
func (h *APIHandler) RatePostUp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pIDstr := vars["post_id"]
	logger.Info.Printf("POST /action/post/%s/up %s\n", pIDstr, r.RemoteAddr)

	session := h.sessionStore.ValidateSession(w, r)

	pID, err := uuid.Parse(pIDstr)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/up - Could not convert id to string: %s\n", pIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pDB, err := h.db.GetPost(h.ctx, pID)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/up - Could not get post: %s\n", pIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.RatePostUp(h.ctx, pDB.ID, session.User.ID)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/up - Could not update post rating: %s\n", pIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info.Printf("OK - POST /action/post/%s/up %s\n", pIDstr, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

// POST /action/post/{Post_GUID}/down
func (h *APIHandler) RatePostDown(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pIDstr := vars["post_id"]
	logger.Info.Printf("POST /action/post/%s/down %s\n", pIDstr, r.RemoteAddr)

	session := h.sessionStore.ValidateSession(w, r)

	pID, err := uuid.Parse(pIDstr)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/down - Could not convert id to string: %s\n", pIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pDB, err := h.db.GetPost(h.ctx, pID)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/down - Could not get post: %s\n", pIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.RatePostDown(h.ctx, pDB.ID, session.User.ID)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/down - Could not update post rating: %s\n", pIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info.Printf("OK - POST /action/post/%s/down %s\n", pIDstr, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}
