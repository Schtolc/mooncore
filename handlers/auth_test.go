package handlers

import (
	"github.com/Schtolc/mooncore/models"
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

var (
	testUser = &models.User{
		Name:     "name5",
		Password: "pass5",
		Email:    "email5@mail.ru",
	}
	token string
)

func TestSignUp(t *testing.T) {
	e := httpexpect.New(t, localhost.String())
	m := map[string]string{
		"Code": "200",
	}
	e.POST("/sign_up").WithJSON(testUser).Expect().Status(http.StatusOK).JSON().Object().Equal(m)

}

func TestSignIn(t *testing.T) {
	e := httpexpect.New(t, localhost.String())
	m := e.POST("/sign_in").WithJSON(testUser).Expect().Status(http.StatusOK).JSON().Object().ContainsKey("Token")
	token = m.Value("Token").String().Raw()
}

func TestSignOut(t *testing.T) {
	e := httpexpect.New(t, localhost.String())
	m := map[string]string{
		"Code": "200",
	}
	e.POST("/sign_out").WithHeader("Authorization", "Bearer "+token).WithJSON(testUser).Expect().Status(http.StatusOK).JSON().Object().Equal(m)
}
