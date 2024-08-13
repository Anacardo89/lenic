package pages

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"math/rand/v2"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Anacardo89/tpsi25_blog/auth"
	"github.com/Anacardo89/tpsi25_blog/internal/model"
	"github.com/Anacardo89/tpsi25_blog/internal/query"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

type PostPage struct {
	Id         int
	GUID       string
	User       string
	Title      string
	RawContent string
	Content    template.HTML
	Image      []byte
	Date       string
	Comments   []Comment
	Session    auth.Session
}

func (p PostPage) TruncatedText() string {
	chars := 0
	for i := range p.RawContent {
		chars++
		if chars > 150 {
			return p.RawContent[:i] + `...`
		}
	}
	return p.RawContent
}

func PostGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p := PostPage{
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

func PostPOST(w http.ResponseWriter, r *http.Request) {
	var fileBytes []byte
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}

	post := model.Post{
		Title:   r.FormValue("post_title"),
		Content: r.FormValue("post_content"),
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		if err != http.ErrMissingFile { // Check for other errors
			logger.Error.Println(err)
			return
		}

		// No image file uploaded
		session := auth.ValidateSession(r)
		post.User = session.User.UserName
		post.GUID = createGUID(post.Title, post.User)

		// Insert post without image data
		_, err = db.Dbase.Exec(query.InsertPost,
			post.GUID, post.Title, post.User, post.Content, []byte{}, "", 1)
		if err != nil {
			logger.Error.Println(err.Error())
			fmt.Fprintln(w, err.Error())
			return
		}

		// Redirect to /home
		http.Redirect(w, r, "/home", http.StatusMovedPermanently)
		return
	}

	// Handle uploaded image
	fileBytes, err = io.ReadAll(file)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	defer file.Close()

	post.Image = fileBytes
	session := auth.ValidateSession(r)
	post.User = session.User.UserName
	post.GUID = createGUID(post.Title, post.User)

	// Insert post with image data
	_, err = db.Dbase.Exec(query.InsertPost,
		post.GUID, post.Title, post.User, post.Content, post.Image, filepath.Ext(header.Filename), 1)
	if err != nil {
		logger.Error.Println(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}

	// Redirect to /home
	http.Redirect(w, r, "/home", http.StatusMovedPermanently)
}

func createGUID(title string, user string) string {
	var guid string
	random := rand.IntN(999)
	guid = strings.ReplaceAll(title, " ", "-")
	guid = guid + strconv.Itoa(random) + user
	return base64.URLEncoding.EncodeToString([]byte(guid))
}
