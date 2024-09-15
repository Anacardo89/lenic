package actions

import (
	"database/sql"
	"encoding/base64"
	"io"
	"math/rand/v2"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/wsoc"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/fsops"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/Anacardo89/tpsi25_blog/pkg/parse"
	"github.com/Anacardo89/tpsi25_blog/pkg/wsocket"
	"github.com/gorilla/mux"
)

func AddPost(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/post ", r.RemoteAddr)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		logger.Error.Println("/action/post - Could not parse Form: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	session := auth.ValidateSession(w, r)

	is_public := false
	visibility := r.FormValue("post-visibility")
	if visibility == "1" {
		is_public = true
	}

	dbpost := database.Post{
		GUID:     createGUID(r.FormValue("post-title"), session.User.UserName),
		Title:    r.FormValue("post-title"),
		Content:  r.FormValue("post-content"),
		AuthorId: session.User.Id,
		Image:    "",
		ImageExt: "",
		IsPublic: is_public,
		Rating:   0,
		Active:   1,
	}
	if dbpost.Title == "" || dbpost.Content == "" {
		http.Error(w, "Post must contain a title and a body", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("post-image")
	if err != nil {
		if err != http.ErrMissingFile {
			logger.Error.Println("/action/post - Could not get image: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = orm.Da.CreatePost(&dbpost)
		if err != nil {
			logger.Error.Println("/action/post - Could not create post: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dbP, err := orm.Da.GetPostByGUID(dbpost.GUID)
		if err != nil {
			logger.Error.Println("/action/post - Could not get post: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dbUser, err := orm.Da.GetUserByID(dbP.AuthorId)
		if err != nil {
			logger.Error.Printf("PUT /action/post/%s/comment - Could not get user: %s\n", dbP.GUID, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mentions := parse.ParseAtString(dbP.Content)
		if len(mentions) > 0 {
			for _, mention := range mentions {
				mention = strings.TrimLeft(mention, "@")
				dbTag, err := orm.Da.GetTagByName(mention)
				if err == sql.ErrNoRows {
					t := &database.Tag{
						TagName: mention,
						TagType: "user",
					}
					err := orm.Da.CreateTag(t)
					if err != nil {
						logger.Error.Printf("PUT /action/post/%s/comment - Could not get create Tag: %s\n", dbP.GUID, err)
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
				} else if err != nil {
					logger.Error.Printf("PUT /action/post/%s/comment - Could not get tag By Id: %s\n", dbP.GUID, err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				ut := &database.UserTag{
					TagId:     dbTag.Id,
					PostId:    dbP.Id,
					CommentId: -1,
					TagPlace:  "post",
				}
				err = orm.Da.CreateUserTag(ut)
				if err != nil {
					logger.Error.Printf("POST /action/post/%s/comment - Could not create UserTag: %s\n", dbP.GUID, err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				wsMsg := wsocket.Message{
					FromUserName: dbUser.UserName,
					Type:         "post_tag",
					Msg:          " has tagged you in their post",
					ResourceId:   dbP.GUID,
					ParentId:     "",
				}

				wsoc.HandlePostTag(wsMsg, mention)
			}
		}

		w.WriteHeader(http.StatusCreated)
		return
	}

	// Handle uploaded image
	dbpost.ImageExt = filepath.Ext(header.Filename)
	fileName := fsops.NameImg(16)
	dbpost.Image = fileName
	imgData, err := io.ReadAll(file)
	if err != nil {
		logger.Error.Println("/action/post - Could not read image data: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fsops.SaveImg(imgData, fsops.PostImgPath, fileName, dbpost.ImageExt)

	// Insert post with image data
	err = orm.Da.CreatePost(&dbpost)
	if err != nil {
		logger.Error.Println("/action/post - Could not not create post: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbP, err := orm.Da.GetPostByGUID(dbpost.GUID)
	if err != nil {
		logger.Error.Println("/action/post - Could not get post: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := orm.Da.GetUserByID(dbP.AuthorId)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment - Could not get user: %s\n", dbP.GUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mentions := parse.ParseAtString(dbP.Content)
	if len(mentions) > 0 {
		for _, mention := range mentions {
			mention = strings.TrimLeft(mention, "@")
			dbTag, err := orm.Da.GetTagByName(mention)
			if err == sql.ErrNoRows {
				t := &database.Tag{
					TagName: mention,
					TagType: "user",
				}
				err := orm.Da.CreateTag(t)
				if err != nil {
					logger.Error.Printf("PUT /action/post/%s/comment - Could not get create Tag: %s\n", dbP.GUID, err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else if err != nil {
				logger.Error.Printf("PUT /action/post/%s/comment - Could not get tag By Id: %s\n", dbP.GUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ut := &database.UserTag{
				TagId:     dbTag.Id,
				PostId:    dbP.Id,
				CommentId: -1,
				TagPlace:  "post",
			}
			err = orm.Da.CreateUserTag(ut)
			if err != nil {
				logger.Error.Printf("POST /action/post/%s/comment - Could not create UserTag: %s\n", dbP.GUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			wsMsg := wsocket.Message{
				FromUserName: dbUser.UserName,
				Type:         "post_tag",
				Msg:          " has tagged you in their post",
				ResourceId:   dbP.GUID,
				ParentId:     "",
			}

			wsoc.HandlePostTag(wsMsg, mention)
		}
	}

	logger.Info.Println("OK - /action/post ", r.RemoteAddr)
	w.WriteHeader(http.StatusCreated)
}

func createGUID(title string, user string) string {
	var guid string
	random := rand.IntN(999)
	guid = strings.ReplaceAll(title, " ", "-")
	guid = guid + strconv.Itoa(random) + user
	return base64.URLEncoding.EncodeToString([]byte(guid))
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	logger.Info.Printf("PUT /action/post/%s %s\n", postGUID, r.RemoteAddr)

	err := r.ParseForm()
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s - Could not parse form: %s\n", postGUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	is_public := false
	visibility := r.FormValue("visibility")
	if visibility == "1" {
		is_public = true
	}

	p := database.Post{
		GUID:     postGUID,
		Title:    r.FormValue("title"),
		Content:  r.FormValue("content"),
		IsPublic: is_public,
	}
	if p.Content == "" || p.Title == "" {
		http.Error(w, "All form fields must be filled out", http.StatusBadRequest)
		return
	}

	err = orm.Da.UpdatePost(p)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s - Could not update post: %s\n", postGUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbP, err := orm.Da.GetPostByGUID(p.GUID)
	if err != nil {
		logger.Error.Println("/action/post - Could not get post: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := orm.Da.GetUserByID(dbP.AuthorId)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment - Could not get user: %s\n", dbP.GUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mentions := parse.ParseAtString(dbP.Content)
	if len(mentions) > 0 {
		for _, mention := range mentions {
			mention = strings.TrimLeft(mention, "@")
			dbTag, err := orm.Da.GetTagByName(mention)
			if err == sql.ErrNoRows {
				t := &database.Tag{
					TagName: mention,
					TagType: "user",
				}
				err := orm.Da.CreateTag(t)
				if err != nil {
					logger.Error.Printf("PUT /action/post/%s/comment - Could not get create Tag: %s\n", dbP.GUID, err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else if err != nil {
				logger.Error.Printf("PUT /action/post/%s/comment - Could not get tag By Id: %s\n", dbP.GUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ut := &database.UserTag{
				TagId:     dbTag.Id,
				PostId:    dbP.Id,
				CommentId: -1,
				TagPlace:  "post",
			}
			err = orm.Da.CreateUserTag(ut)
			if err != nil {
				logger.Error.Printf("POST /action/post/%s/comment - Could not create UserTag: %s\n", dbP.GUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			wsMsg := wsocket.Message{
				FromUserName: dbUser.UserName,
				Type:         "post_tag",
				Msg:          " has tagged you in their post",
				ResourceId:   dbP.GUID,
				ParentId:     "",
			}

			wsoc.HandlePostTag(wsMsg, mention)
		}
	}

	logger.Info.Printf("OK - PUT /action/post/%s %s\n", postGUID, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	logger.Info.Printf("DELETE /action/post/%s %s\n", postGUID, r.RemoteAddr)

	err := orm.Da.DisablePost(postGUID)
	if err != nil {
		logger.Error.Printf("DELETE /action/post/%s - Could not update comment: %s\n", postGUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info.Printf("OK - DELETE /action/post/%s %s\n", postGUID, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func RatePostUp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	logger.Info.Printf("POST /action/post/%s/up %s\n", postGUID, r.RemoteAddr)

	session := auth.ValidateSession(w, r)

	dbpost, err := orm.Da.GetPostByGUID(postGUID)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/up - Could not get post: %s\n", postGUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = orm.Da.RatePostUp(dbpost.Id, session.User.Id)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/up - Could not update post rating: %s\n", postGUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info.Printf("OK - POST /action/post/%s/up %s\n", postGUID, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func RatePostDown(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	logger.Info.Printf("POST /action/post/%s/down %s\n", postGUID, r.RemoteAddr)

	session := auth.ValidateSession(w, r)

	dbpost, err := orm.Da.GetPostByGUID(postGUID)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/down - Could not get post: %s\n", postGUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = orm.Da.RatePostDown(dbpost.Id, session.User.Id)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/down - Could not update post rating: %s\n", postGUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info.Printf("OK - POST /action/post/%s/down %s\n", postGUID, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}
