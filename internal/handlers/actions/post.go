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

	dbpost := database.Post{
		GUID:           createGUID(r.FormValue("post_title"), session.User.UserName),
		Title:          r.FormValue("post_title"),
		Content:        r.FormValue("post_content"),
		User:           session.User.UserName,
		Image:          "",
		ImageExtention: "",
		Active:         1,
	}
	if dbpost.Title == "" || dbpost.Content == "" {
		redirect.RedirectToError(w, r, "Post must contain a title and a body")
		return
	}

	file, header, err := r.FormFile("image")
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
		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}

	// Handle uploaded image
	dbpost.ImageExtention = filepath.Ext(header.Filename)
	fileName := fsops.NameImg(16)
	dbpost.Image = fileName
	imgData, err := io.ReadAll(file)
	if err != nil {
		logger.Error.Println("/action/post - Could not read image data: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	fsops.SaveImg(imgData, fileName, dbpost.ImageExtention)

	// Insert post with image data
	err = orm.Da.CreatePost(&dbpost)
	if err != nil {
		logger.Error.Println("/action/post - Could not not create post: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func createGUID(title string, user string) string {
	var guid string
	random := rand.IntN(999)
	guid = strings.ReplaceAll(title, " ", "-")
	guid = guid + strconv.Itoa(random) + user
	return base64.URLEncoding.EncodeToString([]byte(guid))
}

/*
func PostGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p := presentation.Post{
		GUID: vars["post_guid"],
	}
	err := db.Dbase.QueryRow(query.SelectPostByGUID,
		p.GUID).Scan(
		&p.Title,
		&p.User,
		&p.RawContent,
		&p.Date,
	)
	if err != nil {
		logger.Error.Println(err.Error())
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}
	// p.Content = template.HTML(p.RawContent)
	// TODO or not TODO

}

*/
