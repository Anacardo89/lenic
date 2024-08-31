package pages

import (
	"database/sql"
	"encoding/base64"
	"html/template"
	"net/http"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/internal/model/mapper"
	"github.com/Anacardo89/tpsi25_blog/internal/model/presentation"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/gorilla/mux"
)

type ProfilePage struct {
	Session presentation.Session
	User    presentation.User
	Posts   []presentation.Post
	Follows int
}

func UserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encoded := vars["encoded_user_name"]
	logger.Info.Printf("/user/%s %s\n", encoded, r.RemoteAddr)

	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		logger.Error.Printf("/user/%s - Could not decode user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	userName := string(bytes)
	logger.Info.Printf("/user/%s %s %s\n", encoded, r.RemoteAddr, userName)

	dbuser, err := orm.Da.GetUserByName(userName)
	if err != nil {
		logger.Error.Printf("/user/%s - Could not get user: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	u := mapper.User(dbuser)

	session := auth.ValidateSession(w, r)

	pp := ProfilePage{
		User:    *u,
		Session: session,
	}

	dbfollow, err := orm.Da.GetUserFollows(session.User.Id, u.Id)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error.Printf("/user/%s - Could not get follows: %s\n", encoded, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		} else {
			pp.Follows = -1
		}
	} else if dbfollow != nil { // Check if dbfollow is not nil
		pp.Follows = dbfollow.Status
	} else {
		pp.Follows = -1 // Default or handle case where dbfollow is nil
	}

	var dbposts *[]database.Post
	if (session.User.Id == u.Id) || (dbfollow != nil && dbfollow.Status == 1) {
		dbposts, err = orm.Da.GetUserPosts(u.Id)
		if err != nil {
			logger.Error.Printf("/user/%s - Could not get Posts: %s\n", encoded, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
	} else {
		dbposts, err = orm.Da.GetUserPublicPosts(u.Id)
		if err != nil {
			logger.Error.Printf("/user/%s - Could not get Posts: %s\n", encoded, err)
			redirect.RedirectToError(w, r, err.Error())
			return
		}
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
		pp.Posts = append(pp.Posts, *post)
	}

	t, err := template.ParseFiles("templates/authorized/user-profile.html")
	if err != nil {
		logger.Error.Printf("/user/%s - Could not parse template: %s\n", encoded, err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	t.Execute(w, pp)
}
