package actions

import (
	"encoding/base64"
	"io"
	"math/rand/v2"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/fsops"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

func AddPost(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/post ", r.RemoteAddr)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error.Println("/action/post - Could not parse Form: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	session := auth.ValidateSession(w, r)

	is_public := false
	visibility := r.FormValue("post-visibility")
	if visibility == "1" {
		is_public = true
	}

	dbpost := database.Post{
		GUID:     createGUID(r.FormValue("post-title"), session.User.UserName),
		Title:    r.FormValue("post-title"),
		Content:  r.FormValue("post-content"),
		AuthorId: session.User.Id,
		Image:    "",
		ImageExt: "",
		IsPublic: is_public,
		Rating:   0,
		Active:   1,
	}
	if dbpost.Title == "" || dbpost.Content == "" {
		redirect.RedirectToError(w, r, "Post must contain a title and a body")
		return
	}

	file, header, err := r.FormFile("post-image")
	if err != nil {
		if err != http.ErrMissingFile {
			logger.Error.Println("/action/post - Could not get image: ", err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
		err = orm.Da.CreatePost(&dbpost)
		if err != nil {
			logger.Error.Println("/action/post - Could not create post: ", err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}

	// Handle uploaded image
	dbpost.ImageExt = filepath.Ext(header.Filename)
	fileName := fsops.NameImg(16)
	dbpost.Image = fileName
	imgData, err := io.ReadAll(file)
	if err != nil {
		logger.Error.Println("/action/post - Could not read image data: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fsops.SaveImg(imgData, fsops.PostImgPath, fileName, dbpost.ImageExt)

	// Insert post with image data
	err = orm.Da.CreatePost(&dbpost)
	if err != nil {
		logger.Error.Println("/action/post - Could not not create post: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	logger.Info.Println("OK - /action/post ", r.RemoteAddr)
	w.WriteHeader(http.StatusCreated)
}

func createGUID(title string, user string) string {
	var guid string
	random := rand.IntN(999)
	guid = strings.ReplaceAll(title, " ", "-")
	guid = guid + strconv.Itoa(random) + user
	return base64.URLEncoding.EncodeToString([]byte(guid))
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	logger.Info.Printf("PUT /action/post/%s %s\n", postGUID, r.RemoteAddr)

	err := r.ParseForm()
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s - Could not parse form: %s\n", postGUID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	is_public := false
	visibility := r.FormValue("visibility")
	if visibility == "1" {
		is_public = true
	}

	p := database.Post{
		GUID:     postGUID,
		Title:    r.FormValue("title"),
		Content:  r.FormValue("content"),
		IsPublic: is_public,
	}
	if p.Content == "" || p.Title == "" {
		redirect.RedirectToError(w, r, "All form fields must be filled out")
		return
	}

	err = orm.Da.UpdatePost(p)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s - Could not update post: %s\n", postGUID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	logger.Info.Printf("OK - PUT /action/post/%s %s\n", postGUID, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	logger.Info.Printf("DELETE /action/post/%s %s\n", postGUID, r.RemoteAddr)

	err := orm.Da.DisablePost(postGUID)
	if err != nil {
		logger.Error.Printf("DELETE /action/post/%s - Could not update comment: %s\n", postGUID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	logger.Info.Printf("OK - DELETE /action/post/%s %s\n", postGUID, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func RatePostUp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	logger.Info.Printf("POST /action/post/%s/up %s\n", postGUID, r.RemoteAddr)

	session := auth.ValidateSession(w, r)

	dbpost, err := orm.Da.GetPostByGUID(postGUID)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/up - Could not get post: %s\n", postGUID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	err = orm.Da.RatePostUp(dbpost.Id, session.User.Id)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/up - Could not update post rating: %s\n", postGUID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	logger.Info.Printf("OK - POST /action/post/%s/up %s\n", postGUID, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func RatePostDown(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	logger.Info.Printf("POST /action/post/%s/down %s\n", postGUID, r.RemoteAddr)

	session := auth.ValidateSession(w, r)

	dbpost, err := orm.Da.GetPostByGUID(postGUID)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/down - Could not get post: %s\n", postGUID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	err = orm.Da.RatePostDown(dbpost.Id, session.User.Id)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/down - Could not update post rating: %s\n", postGUID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	logger.Info.Printf("OK - POST /action/post/%s/down %s\n", postGUID, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}
