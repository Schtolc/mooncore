package handlers

import (
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	e := httpexpect.New(t, "http://127.0.0.1:1323")

	resp := Resp{Code: "200", Message: "ECHO_PING"}
	e.GET("/ping").
		Expect().
		Status(http.StatusOK).JSON().Object().Equal(resp)
}
