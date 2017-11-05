package handlers

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestUpload(t *testing.T) {
	e := expect(t)

	content := []byte{45, 34, 2, 54}

	body := e.POST("/upload").WithMultipart().
		WithFileBytes("file", "abc", content).Expect().
		Status(http.StatusOK).JSON().Object()
	body.Value("code").Equal(http.StatusOK)
	body = body.Value("body").Object()

	body.ContainsKey("id").Value("id").NotNull()
	body.ContainsKey("path").Value("path").NotNull()

	path := body.Value("path").String().Raw()

	data, err := ioutil.ReadFile(dependencies.ConfigInstance().Server.UploadStorage + path)
	assert.Nil(t, err)

	assert.Equal(t, data, content)

	os.Remove(dependencies.ConfigInstance().Server.UploadStorage + path)
}
