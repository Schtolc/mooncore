package rest

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/Schtolc/mooncore/utils"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const letterBytes = "0123456789ABCDEF"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randFilename(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// UploadImage uploads image
func UploadImage(c echo.Context) error {
	r := c.Request()

	file, header, err := r.FormFile("file")
	if err != nil {
		logrus.Error(err)
		return utils.SendResponse(c, http.StatusBadRequest, err.Error())
	}
	defer file.Close()

	filename := randFilename(32) + filepath.Ext(header.Filename)

	f, err := os.OpenFile(dependencies.ConfigInstance().Server.UploadStorage+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logrus.Error(err)
		return utils.InternalServerError(c)
	}
	defer f.Close()
	io.Copy(f, file)

	photo := &models.Photo{
		Path: filename,
	}
	if dbc := dependencies.DBInstance().Create(photo); dbc.Error != nil {
		logrus.Println(dbc.Error)
		return utils.InternalServerError(c)
	}

	return utils.SendResponse(c, http.StatusOK, photo)
}
