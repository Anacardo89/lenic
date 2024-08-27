package pages

import (
	"html/template"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

type FeedPage struct {
	Posts   []presentation.Post
	Session presentation.Session
}

func Feed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("/user/%s/feed %s\n", encoded, r.RemoteAddr)
	feed := FeedPage{}
	feed.Session = auth.ValidateSession(w, r)
	dbposts, err := orm.Da.GetPosts()
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
	t, err := template.ParseFiles("templates/feed.html")
	if err != nil {
		logger.Error.Printf("/user/%s/feed - Could not parse template: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, feed)
}
