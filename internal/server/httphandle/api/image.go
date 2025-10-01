package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// /action/image
func (h *APIHandler) GetPostImage(w http.ResponseWriter, r *http.Request) {
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
	pID, err := uuid.Parse(r.URL.Query().Get("post_id"))
	if err != nil {
		fail("parsing post uuid from URL", err, true, http.StatusBadRequest, "invalid path")
		return
	}
	// DB operations
	pDB, err := h.db.GetPost(r.Context(), pID)
	if err != nil {
		fail("dberr: could not get post", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// Early return if no image
	if pDB.PostImage == "" {
		w.WriteHeader(200)
		return
	}
	// Get img
	imgFile, err := h.img.GetImg(true, pDB.PostImage)
	if err != nil {
		fail("failed to get image", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	defer imgFile.Close()
	// Determine MIME type
	ext := filepath.Ext(pDB.PostImage)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)
	// Stream file directly to response
	if _, err := io.Copy(w, imgFile); err != nil {
		fail("failed to write image to response", err, false, http.StatusInternalServerError, "internal error")
		return
	}
}

// /action/profile-pic
func (h *APIHandler) GetProfilePic(w http.ResponseWriter, r *http.Request) {
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
	bytes, err := base64.URLEncoding.DecodeString(r.URL.Query().Get("encoded_username"))
	if err != nil {
		fail("could not decode user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	userName := string(bytes)
	// DB operations
	uDB, err := h.db.GetUserByUserName(h.ctx, userName)
	if err != nil {
		fail("dberr: could not get user", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	if uDB.ProfilePic == "" {
		w.WriteHeader(200)
		return
	}
	// Get img
	imgFile, err := h.img.GetImg(true, uDB.ProfilePic)
	if err != nil {
		fail("failed to get image", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	defer imgFile.Close()
	// Determine MIME type
	ext := filepath.Ext(uDB.ProfilePic)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)
	// Stream file directly to response
	if _, err := io.Copy(w, imgFile); err != nil {
		fail("failed to write image to response", err, false, http.StatusInternalServerError, "internal error")
		return
	}
}

// /action/user/{user_encoded}/profile-pic
func (h *APIHandler) PostProfilePic(w http.ResponseWriter, r *http.Request) {
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
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fail("could not parse form", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	file, header, err := r.FormFile("profile_pic")
	if err != nil {
		fail("could not get image", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	if file == nil || header == nil {
		fail("no file", errors.New("no file"), true, http.StatusBadRequest, "invalid params")
		return
	}
	bytes, err := base64.URLEncoding.DecodeString(vars["encoded_username"])
	if err != nil {
		fail("could not decode user", err, true, http.StatusBadRequest, "invalid user")
		return
	}
	username := string(bytes)
	// DB operations
	uDB, err := h.db.GetUserByUserName(r.Context(), username)
	if err != nil {
		fail("dberr: could not get user", err, true, http.StatusBadRequest, "invalid params")
		return
	}
	// Handle file
	fileExt := filepath.Ext(header.Filename)
	filename := uDB.ID.String()
	filename = fmt.Sprintf("%s%s", filename, fileExt)
	h.img.SaveImg(file, filename)
	h.img.CreatePreview(filename)
	if err := h.db.UpdateProfilePic(r.Context(), username, filename); err != nil {
		fail("dberr: could not update profile pic", err, true, http.StatusInternalServerError, "internal error")
		return
	}
	// Response
	w.WriteHeader(http.StatusOK)
}
