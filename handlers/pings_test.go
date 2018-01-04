package handlers

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	e := expect(t)

	ping := "ECHO_PING"
	e.GET("/ping").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Equal(ping)
}

func TestPingDb(t *testing.T) {
	e := expect(t)

	e.GET("/ping_db").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("data").String().NotEmpty()
}
