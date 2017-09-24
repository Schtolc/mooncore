package handlers

import (
	"testing"
	"github.com/gavv/httpexpect"
	"github.com/Schtolc/mooncore/models"
	"net/http"
)

var (
	testUser  = &models.User {
		Name: "name5",
		Password: "pass5",
		Email: "email5@mail.ru",
	}
	token string
)
func TestSign(t *testing.T) {
	e := httpexpect.New(t, localhost.String())
	e.POST("/v1/sign_up").WithJSON(testUser).Expect().Status(http.StatusOK).Text().Empty()

}

func TestSignIn(t *testing.T) {
	e := httpexpect.New(t, localhost.String())
	m := e.POST("/v1/sign_in").WithJSON(testUser).Expect().Status(http.StatusOK).JSON().Object().ContainsKey("token")
	token = m.Value("token").String().Raw()
}

func TestSignOut (t * testing.T) {
	e := httpexpect.New(t, localhost.String())
	e.POST("/v1/sign_out").WithHeader("Authorization", "Bearer " + token).WithJSON(testUser).Expect().Status(http.StatusOK).Text().Empty()
}


