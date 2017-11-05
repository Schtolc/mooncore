package handlers

import (
	"github.com/labstack/echo"
	"net/http"
	"io"
	"os"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
	"github.com/Schtolc/mooncore/dependencies"
	"path/filepath"
	"github.com/Schtolc/mooncore/models"
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

func UploadImage(c echo.Context) error {
	r := c.Request()

	file, header, err := r.FormFile("file")
	if err != nil {
		logrus.Error(err)
		return sendResponse(c, http.StatusBadRequest, err.Error())
	}
	defer file.Close()

	filename := randFilename(32) + filepath.Ext(header.Filename)

	f, err := os.OpenFile(dependencies.ConfigInstance().Server.UploadStorage + filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logrus.Error(err)
		return internalServerError(c)
	}
	defer f.Close()
	io.Copy(f, file)

	photo := &models.Photo{
		Path: filename,
	}
	if dbc := dependencies.DBInstance().Create(photo); dbc.Error != nil {
		logrus.Println(dbc.Error)
		return internalServerError(c)
	}

	return sendResponse(c, http.StatusOK, photo)
}
