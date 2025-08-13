package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Anacardo89/lenic/internal/db"
	"github.com/Anacardo89/lenic/internal/helpers"
	"github.com/Anacardo89/lenic/internal/server/wshandle"
	"github.com/Anacardo89/lenic/pkg/logger"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Response struct {
	Data string `json:"data"`
}

// POST /action/post/{Post_GUID}/comment
func (h *APIHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDstr := vars["post_id"]
	logger.Info.Printf("POST /action/post/%s/comment %s\n", postIDstr, r.RemoteAddr)
	session := h.sessionStore.ValidateSession(w, r)

	postID, err := uuid.Parse(postIDstr)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment - Could not create comment: %s\n", postID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c := db.Comment{
		PostID:   postID,
		AuthorID: session.User.ID,
		Content:  r.FormValue("comment_text"),
		Rating:   0,
		IsActive: true,
	}
	cID, err := h.db.CreateComment(h.ctx, &c)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment - Could not create comment: %s\n", postID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbComment, err := h.db.GetComment(h.ctx, cID)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment - Could not get comment Id: %s\n", postID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mentions := helpers.ParseAtString(c.Content)
	if len(mentions) > 0 {
		for _, mention := range mentions {
			mention = strings.TrimLeft(mention, "@")
			userDB, err := h.db.GetUserByUserName(h.ctx, mention)
			ut := &db.UserTag{
				UserID:      userDB.ID,
				TargetID:    cID,
				ResourceTpe: db.ResourceComment.String(),
			}
			err = h.db.CreateUserTag(h.ctx, ut)
			if err != nil {
				logger.Error.Printf("POST /action/post/%s/comment - Could not create UserTag: %s\n", postID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			wsMsg := wshandle.Message{
				FromUserName: session.User.UserName,
				Type:         "comment_tag",
				Msg:          " has tagged you in their comment",
				ResourceID:   dbComment.ID.String(),
				ParentID:     postID,
			}

			h.wsHandler.HandleCommentTag(wsMsg, mention)
		}
	}

	resp := Response{
		Data: cID,
	}
	data, err := json.Marshal(&resp)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment - Could not marshal JSON: %s\n", postID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info.Printf("OK - POST /action/post/%s/comment %s\n", postID, r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	logger.Debug.Println(string(data))
	w.Write(data)
}

// PUT /action/post/{Post_GUID}/comment/{comment_id}
func (h *APIHandler) EditComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDstr := vars["post_id"]
	cIDstr := vars["comment_id"]
	logger.Info.Printf("PUT /action/post/%s/comment/%s %s\n", postIDstr, cIDstr, r.RemoteAddr)
	session := h.sessionStore.ValidateSession(w, r)

	cID, err := uuid.Parse(cIDstr)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment/%s - Could not convert id to string: %s\n", postIDstr, cID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	postID, err := uuid.Parse(postIDstr)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment/%s - Could not convert id to string: %s\n", postIDstr, cID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if r.FormValue("comment") == "" {
		http.Error(w, "All form fields must be filled out", http.StatusBadRequest)
		return
	}

	err = h.db.UpdateCommentContent(h.ctx, cID, r.FormValue("comment"))
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment/%s - Could not update comment: %s\n", postIDstr, cID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbComment, err := h.db.GetComment(h.ctx, cID)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment - Could not get comment Id: %s\n", postIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mentions := helpers.ParseAtString(dbComment.Content)
	if len(mentions) > 0 {
		for _, mention := range mentions {
			mention = strings.TrimLeft(mention, "@")
			userDB, err := h.db.GetUserByUserName(h.ctx, mention)
			ut := &db.UserTag{
				UserID:      userDB.ID,
				TargetID:    cID,
				ResourceTpe: db.ResourceComment.String(),
			}
			err = h.db.CreateUserTag(h.ctx, ut)
			if err != nil {
				logger.Error.Printf("PUT /action/post/%s/comment - Could not create UserTag: %s\n", postIDstr, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			wsMsg := wshandle.Message{
				FromUserName: session.User.UserName,
				Type:         "comment_tag",
				Msg:          " has tagged you in their comment",
				ResourceID:   dbComment.ID.String(),
				ParentID:     postID,
			}

			h.wsHandler.HandleCommentTag(wsMsg, mention)
		}
	}

	logger.Info.Printf("OK - PUT /action/post/%s/comment/%s %s\n", postIDstr, cIDstr, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

// DELETE /action/post/{Post_GUID}/comment/{comment_id}
func (h *APIHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDstr := vars["post_id"]
	cIDstr := vars["comment_id"]
	logger.Info.Printf("DELETE /action/post/%s/comment/%s %s\n", postIDstr, cIDstr, r.RemoteAddr)

	cID, err := uuid.Parse(cIDstr)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment/%s - Could not convert id to string: %s\n", postIDstr, cIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.DisableComment(h.ctx, cID)
	if err != nil {
		logger.Error.Printf("DELETE /action/post/%s/comment/%s - Could not update comment: %s\n", postIDstr, cIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbC, err := h.db.GetComment(h.ctx, cID)
	if err != nil {
		logger.Error.Printf("DELETE /action/post/%s/comment/%s - Could not get comment: %s\n", postIDstr, cIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mentions := helpers.ParseAtString(dbC.Content)
	if len(mentions) > 0 {
		for _, mention := range mentions {
			mention = strings.TrimLeft(mention, "@")
			userDB, err := h.db.GetUserByUserName(h.ctx, mention)
			if err != nil {
				logger.Error.Printf("DELETE /action/post/%s/comment - Could not get user By Id: %s\n", postIDstr, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = h.db.DeleteUserTag(h.ctx, userDB.ID, cID)
			if err != nil {
				logger.Error.Printf("DELETE /action/post/%s/comment - Could not delete tag By Id: %s\n", postIDstr, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	logger.Info.Printf("OK - DELETE /action/post/%s/comment/%s %s\n", postIDstr, cIDstr, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

// POST /action/post/{Post_GUID}/comment/{comment_id}/up
func (h *APIHandler) RateCommentUp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDstr := vars["post_id"]
	cIDstr := vars["comment_id"]
	logger.Info.Printf("POST /action/post/%s/comment/%s/up %s\n", postIDstr, cIDstr, r.RemoteAddr)
	session := h.sessionStore.ValidateSession(w, r)

	cID, err := uuid.Parse(cIDstr)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment/%s - Could not convert id to string: %s\n", postIDstr, cID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.RateCommentUp(h.ctx, cID, session.User.ID)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment/%s/up - Could not update comment rating: %s\n", postIDstr, cIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info.Printf("OK - POST /action/post/%s/comment/%s/up %s\n", postIDstr, cIDstr, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

// POST /action/post/{Post_GUID}/comment/{comment_id}/down
func (h *APIHandler) RateCommentDown(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDstr := vars["post_id"]
	cIDstr := vars["comment_id"]
	logger.Info.Printf("POST /action/post/%s/comment/%s/down %s\n", postIDstr, cIDstr, r.RemoteAddr)
	session := h.sessionStore.ValidateSession(w, r)

	cID, err := uuid.Parse(cIDstr)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment/%s - Could not convert id to string: %s\n", postIDstr, cID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.RateCommentDown(h.ctx, cID, session.User.ID)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment/%s/down - Could not update comment rating: %s\n", postIDstr, cIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info.Printf("OK - POST /action/post/%s/comment/%s/down %s\n", postIDstr, cIDstr, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}
