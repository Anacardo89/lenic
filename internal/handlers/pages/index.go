package pages

import (
	"html/template"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/auth"
	"github.com/Anacardo89/tpsi25_blog/internal/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/actions"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

type IndexPage struct {
	Posts   []actions.PostPage
	Session auth.Session
}

func Index(w http.ResponseWriter, r *http.Request) {
	index := IndexPage{}
	index.Session = auth.ValidateSession(r)
	rows, err := db.Dbase.Query(query.SelectPosts)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		thisPost := actions.PostPage{}
		err := rows.Scan(
			&thisPost.GUID,
			&thisPost.Title,
			&thisPost.User,
			&thisPost.RawContent,
			&thisPost.Date,
		)
		if err != nil {
			logger.Error.Println(err)
			return
		}
		thisPost.Content = template.HTML(thisPost.RawContent)
		index.Posts = append(index.Posts, thisPost)
	}
	t, err := template.ParseFiles("../templates/index.html")
	if err != nil {
		logger.Error.Println(err)
	}
	t.Execute(w, index)
}
