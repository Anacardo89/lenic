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
)

type IndexPage struct {
	Posts   []presentation.Post
	Session presentation.Session
}

func Index(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/index ", r.RemoteAddr)
	index := IndexPage{}
	index.Session = auth.ValidateSession(w, r)
	dbposts, err := orm.Da.GetPosts()
	if err != nil {
		logger.Error.Println("/index - Error getting Posts: ", err)
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
		post := mapper.Post(&dbpost, dbuser.UserName)
		post.Content = template.HTML(post.RawContent)
		index.Posts = append(index.Posts, *post)
	}
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		logger.Error.Println("/index - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, index)
}
