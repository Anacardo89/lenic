package pages

import (
	"html/template"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
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
	index := IndexPage{}
	index.Session = auth.ValidateSession(r)
	dbposts, err := orm.Da.GetPosts()
	if err != nil {
		logger.Error.Println(err)
		return
	}
	for _, dbpost := range *dbposts {
		post := mapper.Post(&dbpost)
		post.Content = template.HTML(post.RawContent)
		index.Posts = append(index.Posts, *post)
	}
	t, err := template.ParseFiles("../templates/index.html")
	if err != nil {
		logger.Error.Println(err)
		return
	}
	t.Execute(w, index)
}
