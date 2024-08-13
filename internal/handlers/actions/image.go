package actions

import (
	"database/sql"
	"mime"
	"net/http"
	"strings"

	"github.com/Anacardo89/tpsi25_blog/pkg/db"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

func Image(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		return
	}

	var postImage []byte
	var imageExtension string
	err := db.Dbase.QueryRow("SELECT post_image, post_image_ext FROM posts WHERE post_guid = ?", guid).Scan(&postImage, &imageExtension)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		logger.Error.Println(err)
		return
	}
	imageExtension = strings.TrimPrefix(imageExtension, ".")
	mimeType := mime.TypeByExtension(imageExtension)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)
	w.Write(postImage)
}
