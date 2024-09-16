package actions

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/wsoc"
	"github.com/Anacardo89/tpsi25_blog/internal/model/database"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
	"github.com/Anacardo89/tpsi25_blog/pkg/parse"
	"github.com/Anacardo89/tpsi25_blog/pkg/wsocket"
	"github.com/gorilla/mux"
)

type Response struct {
	Data string `json:"data"`
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	logger.Info.Printf("POST /action/post/%s/comment %s\n", postGUID, r.RemoteAddr)
	session := auth.ValidateSession(w, r)

	c := database.Comment{
		PostGUID: postGUID,
		AuthorId: session.User.Id,
		Content:  r.FormValue("comment_text"),
		Rating:   0,
		Active:   1,
	}
	res, err := orm.Da.CreateComment(&c)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment - Could not create comment: %s\n", postGUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment - Could not get notification Id: %s\n", postGUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbComment, err := orm.Da.GetCommentById(int(lastInsertID))
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment - Could not get comment Id: %s\n", postGUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mentions := parse.ParseAtString(c.Content)
	if len(mentions) > 0 {
		for _, mention := range mentions {
			mention = strings.TrimLeft(mention, "@")
			_, err := orm.Da.GetTagByName(mention)
			if err == sql.ErrNoRows {
				t := &database.Tag{
					TagName: mention,
					TagType: "user",
				}
				err := orm.Da.CreateTag(t)
				if err != nil {
					logger.Error.Printf("POST /action/post/%s/comment - Could not get create Tag: %s\n", postGUID, err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else if err != nil {
				logger.Error.Printf("POST /action/post/%s/comment - Could not get tag By Id: %s\n", postGUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			dbTag, err := orm.Da.GetTagByName(mention)
			if err != nil {
				logger.Error.Printf("POST /action/post/%s/comment - Could not get create Tag: %s\n", postGUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			dbPost, err := orm.Da.GetPostByGUID(dbComment.PostGUID)
			if err != nil {
				logger.Error.Printf("POST /action/post/%s/comment - Could not get Post by Id: %s\n", postGUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ut := &database.UserTag{
				TagId:     dbTag.Id,
				PostId:    dbPost.Id,
				CommentId: dbComment.Id,
				TagPlace:  "comment",
			}
			err = orm.Da.CreateUserTag(ut)
			if err != nil {
				logger.Error.Printf("POST /action/post/%s/comment - Could not create UserTag: %s\n", postGUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			idString := strconv.Itoa(dbComment.Id)
			wsMsg := wsocket.Message{
				FromUserName: session.User.UserName,
				Type:         "comment_tag",
				Msg:          " has tagged you in their comment",
				ResourceId:   idString,
				ParentId:     dbPost.GUID,
			}

			wsoc.HandleCommentTag(wsMsg, mention)
		}
	}

	idstring := strconv.Itoa(int(lastInsertID))
	resp := Response{
		Data: idstring,
	}
	data, err := json.Marshal(&resp)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment - Could not marshal JSON: %s\n", postGUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Info.Printf("OK - POST /action/post/%s/comment %s\n", postGUID, r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	logger.Debug.Println(string(data))
	w.Write(data)
}

func EditComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	id := vars["comment_id"]
	logger.Info.Printf("PUT /action/post/%s/comment/%s %s\n", postGUID, id, r.RemoteAddr)

	idint, err := strconv.Atoi(id)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment/%s - Could not convert id to string: %s\n", postGUID, id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = r.ParseForm()
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment/%s - Could not parse form: %s\n", postGUID, id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c := database.Comment{
		Id:      idint,
		Content: r.FormValue("comment"),
	}
	if c.Content == "" {
		http.Error(w, "All form fields must be filled out", http.StatusBadRequest)
		return
	}

	err = orm.Da.UpdateCommentText(c.Id, c.Content)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment/%s - Could not update comment: %s\n", postGUID, id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbComment, err := orm.Da.GetCommentById(c.Id)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment - Could not get comment Id: %s\n", postGUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := orm.Da.GetUserByID(dbComment.AuthorId)
	if err != nil {
		logger.Error.Printf("PUT /action/post/%s/comment - Could not get user: %s\n", postGUID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mentions := parse.ParseAtString(c.Content)
	if len(mentions) > 0 {
		for _, mention := range mentions {
			mention = strings.TrimLeft(mention, "@")
			_, err := orm.Da.GetTagByName(mention)
			if err == sql.ErrNoRows {
				t := &database.Tag{
					TagName: mention,
					TagType: "user",
				}
				err := orm.Da.CreateTag(t)
				if err != nil {
					logger.Error.Printf("PUT /action/post/%s/comment - Could not get create Tag: %s\n", postGUID, err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			} else if err != nil {
				logger.Error.Printf("PUT /action/post/%s/comment - Could not get tag By Id: %s\n", postGUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			dbTag, err := orm.Da.GetTagByName(mention)
			if err != nil {
				logger.Error.Printf("POST /action/post/%s/comment - Could not get create Tag: %s\n", postGUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			dbPost, err := orm.Da.GetPostByGUID(dbComment.PostGUID)
			if err != nil {
				logger.Error.Printf("PUT /action/post/%s/comment - Could not get Post by Id: %s\n", postGUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ut := &database.UserTag{
				TagId:     dbTag.Id,
				PostId:    dbPost.Id,
				CommentId: dbComment.Id,
				TagPlace:  "comment",
			}
			err = orm.Da.CreateUserTag(ut)
			if err != nil {
				logger.Error.Printf("POST /action/post/%s/comment - Could not create UserTag: %s\n", postGUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			idString := strconv.Itoa(dbComment.Id)
			wsMsg := wsocket.Message{
				FromUserName: dbUser.UserName,
				Type:         "comment_tag",
				Msg:          " has tagged you in their comment",
				ResourceId:   idString,
				ParentId:     dbPost.GUID,
			}

			wsoc.HandleCommentTag(wsMsg, mention)
		}
	}

	logger.Info.Printf("OK - PUT /action/post/%s/comment/%s %s\n", postGUID, id, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	id := vars["comment_id"]
	logger.Info.Printf("DELETE /action/post/%s/comment/%s %s\n", postGUID, id, r.RemoteAddr)

	idint, err := strconv.Atoi(id)
	if err != nil {
		logger.Error.Printf("DELETE /action/post/%s/comment/%s - Could not convert id to string: %s\n", postGUID, id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = orm.Da.DisableComment(idint)
	if err != nil {
		logger.Error.Printf("DELETE /action/post/%s/comment/%s - Could not update comment: %s\n", postGUID, id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbC, err := orm.Da.GetCommentById(idint)
	if err != nil {
		logger.Error.Printf("DELETE /action/post/%s/comment/%s - Could not get comment: %s\n", postGUID, id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mentions := parse.ParseAtString(dbC.Content)
	if len(mentions) > 0 {
		for _, mention := range mentions {
			mention = strings.TrimLeft(mention, "@")
			dbTag, err := orm.Da.GetTagByName(mention)
			if err != nil {
				logger.Error.Printf("DELETE /action/post/%s/comment - Could not get tag By Id: %s\n", postGUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = orm.Da.DeleteUserTagByID(dbTag.Id)
			if err != nil {
				logger.Error.Printf("DELETE /action/post/%s/comment - Could not delete tag By Id: %s\n", postGUID, err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	logger.Info.Printf("OK - DELETE /action/post/%s/comment/%s %s\n", postGUID, id, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func RateCommentUp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	id := vars["comment_id"]
	logger.Info.Printf("POST /action/post/%s/comment/%s/up %s\n", postGUID, id, r.RemoteAddr)
	comment_id, err := strconv.Atoi(id)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment/%s/up - Could not convert id to string: %s\n", postGUID, id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session := auth.ValidateSession(w, r)
	err = orm.Da.RateCommentUp(comment_id, session.User.Id)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment/%s/up - Could not update comment rating: %s\n", postGUID, id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info.Printf("OK - POST /action/post/%s/comment/%s/up %s\n", postGUID, id, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}

func RateCommentDown(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postGUID := vars["post_guid"]
	id := vars["comment_id"]
	logger.Info.Printf("POST /action/post/%s/comment/%s/down %s\n", postGUID, id, r.RemoteAddr)
	comment_id, err := strconv.Atoi(id)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment/%s/down - Could not convert id to string: %s\n", postGUID, id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session := auth.ValidateSession(w, r)
	err = orm.Da.RateCommentDown(comment_id, session.User.Id)
	if err != nil {
		logger.Error.Printf("POST /action/post/%s/comment/%s/down - Could not update comment rating: %s\n", postGUID, id, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Info.Printf("OK - POST /action/post/%s/comment/%s/down %s\n", postGUID, id, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
}
