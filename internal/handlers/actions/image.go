package actions

import (
	"database/sql"
	"mime"
	"net/http"
	"strings"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

func Image(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		return
	}
	dbpost, err := orm.Da.GetPostByGUID(guid)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		logger.Error.Println(err)
		return
	}
	dbpost.ImageExtention = strings.TrimPrefix(dbpost.ImageExtention, ".")
	mimeType := mime.TypeByExtension(dbpost.ImageExtention)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)
	w.Write(dbpost.Image)
}
