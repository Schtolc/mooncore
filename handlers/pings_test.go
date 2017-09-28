package handlers

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	e := expect(t)

	ping := &Response{ Code:  OK, Body: "ECHO_PING" }
	e.GET("/ping").
		Expect().
		Status(http.StatusOK).JSON().Object().Equal(ping)
}

func TestPingDb(t *testing.T) {
	e := expect(t)

	e.GET("/ping_db").
		Expect().
		Status(http.StatusOK).
		JSON().Object().Value("body").String().NotEmpty()
}
