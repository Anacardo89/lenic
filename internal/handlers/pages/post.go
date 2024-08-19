package pages

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

func NewPost(w http.ResponseWriter, r *http.Request) {
	post := presentation.Post{
		Session: auth.ValidateSession(w, r),
	}
	t, err := template.ParseFiles("templates/newPost.html")
	if err != nil {
		logger.Error.Println(err)
		return
	}
	t.Execute(w, post)
}

func Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	dbpost, err := orm.Da.GetPostByGUID(postGUID)
	if err != nil {
		logger.Error.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}
	p := mapper.Post(dbpost)
	p.Session = auth.ValidateSession(w, r)
	p.Content = template.HTML(p.RawContent)

	dbcomments, err := orm.Da.GetCommentsByPost(p.GUID)
	if err != nil {
		logger.Error.Println(err)
	}
	for _, dbcomment := range *dbcomments {
		c := mapper.Comment(&dbcomment)
		p.Comments = append(p.Comments, *c)
	}
	t, err := template.ParseFiles("templates/post.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	t.Execute(w, p)
}
