package test

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"math/rand"
	"github.com/Schtolc/mooncore/dao"
)

func TestUpload(t *testing.T) {
	e := expect(t)

	length := rand.Intn(50) + 10

	content := make([]byte, length)

	for i := range content {
		content[i] = byte(rand.Int())
	}

	body := e.POST("/upload").WithMultipart().
		WithFileBytes("file", "abc", content).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object()

	body.ContainsKey("id").Value("id").NotNull()
	body.ContainsKey("path").Value("path").NotNull()

	id := int64(body.Value("id").Number().Raw())
	path := body.Value("path").String().Raw()

	data, err := ioutil.ReadFile(dependencies.ConfigInstance().Server.UploadStorage + path)

	assert.Nil(t, err, "cannot read file")
	assert.Equal(t, data, content, "content differs")

	if err := dao.DeletePhoto(id); err != nil {
		t.Error("cannot delete photo")
	}

	if err := os.Remove(dependencies.ConfigInstance().Server.UploadStorage + path); err != nil {
		t.Error("cannot delete file")
	}
}
