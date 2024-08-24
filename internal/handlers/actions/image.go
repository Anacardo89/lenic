package actions

import (
	"database/sql"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/Anacardo89/tpsi25_blog/internal/handlers/data/orm"
	"github.com/Anacardo89/tpsi25_blog/internal/handlers/redirect"
	"github.com/Anacardo89/tpsi25_blog/pkg/fsops"
	"github.com/Anacardo89/tpsi25_blog/pkg/logger"
)

func Image(w http.ResponseWriter, r *http.Request) {
	logger.Info.Println("/action/image ", r.RemoteAddr)
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		return
	}
	dbpost, err := orm.Da.GetPostByGUID(guid)
	if err != nil {
		if err == sql.ErrNoRows {
			return
		}
		logger.Error.Println("/action/image - Could not get post: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	imgPath := fsops.ImgPath + dbpost.Image + dbpost.ImageExt
	imgFile, err := os.Open(imgPath)
	if err != nil {
		logger.Error.Println("/action/image - Could not open image: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}
	defer imgFile.Close()

	imgData, err := io.ReadAll(imgFile)
	if err != nil {
		logger.Error.Println("/action/image - Could not read image: ", err)
		redirect.RedirectToError(w, r, err.Error())
		return
	}

	dbpost.ImageExt = strings.TrimPrefix(dbpost.ImageExt, ".")
	mimeType := mime.TypeByExtension(dbpost.ImageExt)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)
	w.Write(imgData)
}
