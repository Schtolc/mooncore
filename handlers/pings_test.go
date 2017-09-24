package handlers

import (
	"github.com/Schtolc/mooncore/utils"
	"github.com/gavv/httpexpect"
	"net/http"
	"net/url"
	"testing"
)

var (
	conf      = utils.GetConfig()
	localhost = url.URL{Scheme: "http", Host: conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port}
)

func TestPing(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	e.GET("/v1/ping").
		Expect().
		Status(http.StatusOK).Text().Equal("ECHO_PING")
}

func TestPingDb(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	e.GET("/v1/ping_db").
		Expect().
		Status(http.StatusOK).JSON().Object().ContainsKey("message").NotEmpty()
}
