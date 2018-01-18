package test

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"net/http"
	"testing"
)

var (
	authUser = &models.UserAuth{
		Name:     "name5",
		Password: "pass5",
		Email:    "email5@mail.ru",
	}
	Stranger = &models.UserAuth{
		Name:     "1",
		Password: "1",
		Email:    "1@mail.ru",
	}
	token string
)

func TestSignUp(t *testing.T) {
	e := expect(t)

	e.POST("/sign_up").WithJSON(authUser).Expect().Status(http.StatusOK)
	e.POST("/sign_up").WithJSON(authUser).Expect().Status(http.StatusBadRequest)
	e.POST("/sign_up").WithText("unexpected request body").Expect().Status(http.StatusBadRequest)

	e.POST("/sign_up").WithJSON(&map[string]string{
		"qwe": "qwe",
	}).Expect().Status(http.StatusBadRequest)
}

func TestSignIn(t *testing.T) {
	e := expect(t)

	m := e.POST("/sign_in").WithJSON(authUser).Expect().Status(http.StatusOK).JSON().Object().ContainsKey("data")
	token = m.Value("data").String().Raw()

	e.POST("/sign_in").WithJSON(Stranger).Expect().Status(http.StatusBadRequest)
}

func TestPingAuth(t *testing.T) {
	e := expect(t)

	e.POST("/auth_ping").WithHeader("Authorization", "Bearer "+token).WithJSON(authUser).
		Expect().Status(http.StatusOK).JSON().Object().Value("data").Equal(authUser.Name)

	dependencies.DBInstance().Delete(&authUser)
}
