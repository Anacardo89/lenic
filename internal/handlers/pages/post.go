package pages

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/auth"
	"github.com/Anacardo89/tpsi25_blog/internal/data/query"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/actions"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

func NewPost(w http.ResponseWriter, r *http.Request) {
	postpg := actions.PostPage{
		Session: auth.ValidateSession(r),
	}
	t, err := template.ParseFiles("../templates/newPost.html")
	if err != nil {
		logger.Error.Println(err)
	}
	t.Execute(w, postpg)
}

func Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p := actions.PostPage{
		Session: auth.ValidateSession(r),
		GUID:    vars["post_guid"],
	}
	err := db.Dbase.QueryRow(query.SelectPostByGUID, p.GUID).Scan(
		&p.Title,
		&p.User,
		&p.RawContent,
		&p.Image,
		&p.Date,
	)
	if err != nil {
		logger.Error.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}
	p.Content = template.HTML(p.RawContent)
	comments, err := db.Dbase.Query(query.SelectComments, p.GUID)
	if err != nil {
		logger.Error.Println(err)
	}
	for comments.Next() {
		var c actions.Comment
		comments.Scan(
			&c.Id,
			&c.UserName,
			&c.CommentText,
			&c.Date,
		)
		p.Comments = append(p.Comments, c)
	}
	t, err := template.ParseFiles("../templates/post.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	t.Execute(w, p)

}
