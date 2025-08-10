package page

import (
	"encoding/base64"
	"html/template"
	"net/http"

	"github.com/Anacardo89/lenic/internal/handlers/data/orm"
	"github.com/Anacardo89/lenic/internal/handlers/redirect"
	"github.com/Anacardo89/lenic/internal/model/mapper"
	"github.com/Anacardo89/lenic/internal/model/presentation"
	"github.com/Anacardo89/lenic/pkg/auth"
	"github.com/Anacardo89/lenic/pkg/logger"
	"github.com/gorilla/mux"
)

type FeedPage struct {
	Session presentation.Session
	Posts   []presentation.Post
}

func (h *PageHandler) Feed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("/user/%s/feed %s\n", encoded, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("/user/%s/feed - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("/user/%s/feed %s %s\n", encoded, r.RemoteAddr, userName)

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("/user/%s/feed - Could not get user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	feed := FeedPage{}
	feed.Session = auth.ValidateSession(w, r)
	dbposts, err := orm.Da.GetFeed(dbuser.Id)
	if err != nil {
		logger.Error.Printf("/user/%s/feed - Could not get Posts: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	for _, dbpost := range *dbposts {
		dbuser, err := orm.Da.GetUserByID(dbpost.AuthorId)
		if err != nil {
			logger.Error.Printf("/post/%s - Could not get Comment Author: %s\n", dbpost.GUID, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
		u := mapper.User(dbuser)
		post := mapper.Post(&dbpost, u)
		post.Content = template.HTML(post.RawContent)
		feed.Posts = append(feed.Posts, *post)
	}
	t, err := template.ParseFiles("templates/authorized/feed.html")
	if err != nil {
		logger.Error.Printf("/user/%s/feed - Could not parse template: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, feed)
}
