package actions

import (
	"database/sql"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/pkg/fsops"
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

	imgPath := fsops.ImgPath + dbpost.Image + dbpost.ImageExtention
	imgFile, err := os.Open(imgPath)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	defer imgFile.Close()

	imgData, err := io.ReadAll(imgFile)
	if err != nil {
		logger.Error.Println(err)
	}

	dbpost.ImageExtention = strings.TrimPrefix(dbpost.ImageExtention, ".")
	mimeType := mime.TypeByExtension(dbpost.ImageExtention)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)
	w.Write(imgData)
}
