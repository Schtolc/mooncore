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

	resp := Resp{Code: "200", Message: "ECHO_PING"}
	e.GET("/ping").
		Expect().
		Status(http.StatusOK).JSON().Object().Equal(resp)
}

func TestPingDb(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	resp := Resp{Code: "200", Message: "1"}
	e.GET("/ping_db").
		Expect().
		Status(http.StatusOK).JSON().Object().Equal(resp)
}
