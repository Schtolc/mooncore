package handlers

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	e := expect(t)

	resp := Resp{Code: "200", Message: "ECHO_PING"}
	e.GET("/ping").
		Expect().
		Status(http.StatusOK).
		JSON().Object().Equal(resp)
}

func TestPingDb(t *testing.T) {
	e := expect(t)

	e.GET("/ping_db").
		Expect().
		Status(http.StatusOK).
		JSON().Object().Value("message").String().NotEmpty()
}
