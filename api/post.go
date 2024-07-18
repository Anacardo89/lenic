package api

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

	"github.com/Anacardo89/tpsi25_blog.git/auth"
	"github.com/Anacardo89/tpsi25_blog.git/db"
	"github.com/Anacardo89/tpsi25_blog.git/logger"
	"github.com/gorilla/mux"
)

type PostPage struct {
	Id         int
	GUID       string
	User       string
	Title      string
	RawContent string
	Content    template.HTML
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
	err := db.Dbase.QueryRow(db.SelectPostByGUID,
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
	p.Content = template.HTML(p.RawContent)

	// TODO or not TODO

}

func PostPOST(w http.ResponseWriter, r *http.Request) {
	var fileBytes []byte
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	post := db.Post{
		PostTitle:   r.FormValue("post_title"),
		PostContent: r.FormValue("post_content"),
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
		} else {
			logger.Error.Println(err)
			return
		}
	} else {
		fileBytes, err = io.ReadAll(file)
		if err != nil {
			logger.Error.Println(err)
			return
		}
	}
	defer file.Close()
	post.PostImage = fileBytes
	session := auth.ValidateSession(r)
	post.PostUser = session.User.UserName
	post.PostGUID = createGUID(post.PostTitle, post.PostUser)
	_, err = db.Dbase.Exec(db.InsertPost,
		post.PostGUID, post.PostTitle, post.PostUser, post.PostContent, post.PostImage, filepath.Ext(header.Filename), 1)
	if err != nil {
		logger.Error.Println(err.Error())
		fmt.Fprintln(w, err.Error())
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
