package api

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Anacardo89/lenic/pkg/fsops"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// /action/image
func (h *APIHandler) PostImage(w http.ResponseWriter, r *http.Request) {
	pIDstr := r.URL.Query().Get("post_id")
	if pIDstr == "" {
		return
	}

	pID, err := uuid.Parse(pIDstr)
	if err != nil {
		logger.Error.Printf("/action/image - Could not convert id to string: %s\n", pIDstr, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pDB, err := h.db.GetPost(h.ctx, pID)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		logger.Error.Println("/action/image - Could not get post: ", err)
		return
	}

	imgPath := fsops.PostImgPath + pDB.PostImage
	imgFile, err := os.Open(imgPath)
	if err != nil {
		logger.Error.Println("/action/image - Could not open image: ", err)
		return
	}
	defer imgFile.Close()

	imgData, err := io.ReadAll(imgFile)
	if err != nil {
		logger.Error.Println("/action/image - Could not read image: ", err)
		return
	}

	imageExt := strings.TrimPrefix(pDB.PostImage, ".")
	mimeType := mime.TypeByExtension(imageExt)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)
	w.Write(imgData)
}

// /action/profile-pic
func (h *APIHandler) ProfilePic(w http.ResponseWriter, r *http.Request) {
	encoded := r.URL.Query().Get("encoded_username")
	if encoded == "" {
		return
	}

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("/action/profile-pic - Could not decode user: %s\n", err)
		return
	}
	userName := string(bytes)

	uDB, err := h.db.GetUserByUserName(h.ctx, userName)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		logger.Error.Println("/action/profile-pic - Could not get user: ", err)
		return
	}

	imgPath := fsops.ProfilePicPath + uDB.ProfilePic
	imgFile, err := os.Open(imgPath)
	if err != nil {
		logger.Error.Println("/action/profile-pic - Could not open image: ", err)
		return
	}
	defer imgFile.Close()

	imgData, err := io.ReadAll(imgFile)
	if err != nil {
		logger.Error.Println("/action/image - Could not read image: ", err)
		return
	}

	picExt := strings.TrimPrefix(uDB.ProfilePic, ".")
	mimeType := mime.TypeByExtension(picExt)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)
	w.Write(imgData)
}

// /action/user/{user_encoded}/profile-pic
func (h *APIHandler) PostProfilePic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_username"]

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error.Printf("/action/user/%s/profile-pic - Could not parse form  %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("profile_pic")
	if err != nil {
		logger.Error.Printf("/action/user/%s/profile-pic - Could not get image: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileName := fsops.NameImg(16)
	fileExt := filepath.Ext(header.Filename)
	fileName = fmt.Sprintf("%s.%s", fileName, fileExt)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("/action/user/%s/profile-pic - Could not decode user: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userName := string(bytes)

	err = h.db.UpdateProfilePic(h.ctx, userName, fileName)
	if err != nil {
		logger.Error.Printf("/action/user/%s/profile-pic - Could not update profile pic: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	imgData, err := io.ReadAll(file)
	if err != nil {
		logger.Error.Println("/action/image - Could not read image: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fsops.SaveImg(imgData, fsops.ProfilePicPath, fileName)
	w.WriteHeader(http.StatusOK)
}
