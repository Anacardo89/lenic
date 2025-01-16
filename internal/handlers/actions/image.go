package actions

import (
	"database/sql"
	"encoding/base64"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/pkg/fsops"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

// /action/image
func PostImage(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/image ", r.RemoteAddr)
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		return
	}
	dbpost, err := orm.Da.GetPostByGUID(guid)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		logger.Error.Println("/action/image - Could not get post: ", err)
		return
	}

	imgPath := fsops.PostImgPath + dbpost.Image + dbpost.ImageExt
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

	dbpost.ImageExt = strings.TrimPrefix(dbpost.ImageExt, ".")
	mimeType := mime.TypeByExtension(dbpost.ImageExt)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)
	logger.Info.Println("OK - /action/image ", r.RemoteAddr)
	w.Write(imgData)
}

// /action/profile-pic
func ProfilePic(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/profile-pic ", r.RemoteAddr)
	encoded := r.URL.Query().Get("user-encoded")
	if encoded == "" {
		return
	}

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("/action/profile-pic - Could not decode user: %s\n", err)
		return
	}
	userName := string(bytes)

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		logger.Error.Println("/action/profile-pic - Could not get user: ", err)
		return
	}

	imgPath := fsops.ProfilePicPath + dbuser.ProfilePic + dbuser.ProfilePicExt
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

	dbuser.ProfilePicExt = strings.TrimPrefix(dbuser.ProfilePicExt, ".")
	mimeType := mime.TypeByExtension(dbuser.ProfilePicExt)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)
	logger.Info.Println("OK - /action/profile-pic ", r.RemoteAddr)
	w.Write(imgData)
}

// /action/user/{user_encoded}/profile-pic
func PostProfilePic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("/action/user/%s/profile-pic  %s\n", encoded, r.RemoteAddr)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error.Printf("/action/user/%s/profile-pic - Could not parse form  %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("profile-image")
	if err != nil {
		logger.Error.Printf("/action/user/%s/profile-pic - Could not get image: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileName := fsops.NameImg(16)
	fileExt := filepath.Ext(header.Filename)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("/action/user/%s/profile-pic - Could not decode user: %s\n", encoded, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userName := string(bytes)

	err = orm.Da.UpdateProfilePic(fileName, fileExt, userName)
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
	fsops.SaveImg(imgData, fsops.ProfilePicPath, fileName, fileExt)
	logger.Info.Printf("/action/user/%s/profile-pic  %s\n", encoded, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}
