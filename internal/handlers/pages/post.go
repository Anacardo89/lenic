package pages

import (
	"database/sql"
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

type PostPage struct {
	Session presentation.Session
	Post    presentation.Post
}

func NewPost(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/newPost ", r.RemoteAddr)
	postp := PostPage{
		Session: auth.ValidateSession(w, r),
		Post:    presentation.Post{},
	}
	t, err := template.ParseFiles("templates/authorized/newPost.html")
	if err != nil {
		logger.Error.Println("/newPost - Could not parse template: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, postp)
}

func Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	logger.Info.Printf("/post/%s %s\n", postGUID, r.RemoteAddr)
	postp := PostPage{}
	dbpost, err := orm.Da.GetPostByGUID(postGUID)
	if err != nil {
		logger.Error.Printf("/post/%s - Could not get Post: %s\n", postGUID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	dbuser, err := orm.Da.GetUserByID(dbpost.AuthorId)
	if err != nil {
		logger.Error.Printf("/post/%s - Could not get Author: %s\n", postGUID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	u := mapper.User(dbuser)

	p := mapper.Post(dbpost, u)

	postp.Session = auth.ValidateSession(w, r)
	pr, err := orm.Da.GetPostUserRating(dbpost.Id, postp.Session.User.Id)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error.Printf("/post/%s - Could not get rating: %s\n", postGUID, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
	}
	if pr != nil {
		p.UserRating = pr.RatingValue
	} else {
		p.UserRating = 0
	}
	p.Content = template.HTML(p.RawContent)
	p.Comments = []presentation.Comment{}

	dbcomments, err := orm.Da.GetCommentsByPost(p.GUID)
	if err != nil {
		logger.Error.Println(err)
	}
	for _, dbcomment := range *dbcomments {
		dbuser, err := orm.Da.GetUserByID(dbcomment.AuthorId)
		if err != nil {
			logger.Error.Printf("/post/%s - Could not get Comment Author: %s\n", postGUID, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
		u := mapper.User(dbuser)

		c := mapper.Comment(&dbcomment, u)

		cr, err := orm.Da.GetCommentUserRating(dbcomment.Id, postp.Session.User.Id)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.Error.Printf("/post/%s - Could not get comment rating: %s\n", postGUID, err)
				redirect.RedirectToError(w, r, err.Error())
				return
			}
		}

		if cr != nil {
			c.UserRating = cr.RatingValue
		} else {
			c.UserRating = 0
		}
		p.Comments = append(p.Comments, *c)
	}
	postp.Post = *p
	t, err := template.ParseFiles("templates/authorized/post.html")
	if err != nil {
		logger.Error.Printf("/post/%s - Could not parse template: %s\n", postGUID, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, postp)
}
