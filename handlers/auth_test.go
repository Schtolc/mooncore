package handlers

import (
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"net/http"
	"testing"
)

var (
	testUser = &models.UserAuth{
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

	e.POST("/sign_up").WithJSON(testUser).Expect().Status(http.StatusOK).JSON().Object().Value("code").Equal(http.StatusOK)
	e.POST("/sign_up").WithJSON(testUser).Expect().Status(http.StatusOK).JSON().Object().Value("code").Equal(http.StatusBadRequest)
	e.POST("/sign_up").WithText("unexpected request body").Expect().Status(http.StatusOK).JSON().Object().Value("code").Equal(http.StatusBadRequest)

	e.POST("/sign_up").WithJSON(&map[string]string{
		"qwe": "qwe",
	}).Expect().Status(http.StatusOK).JSON().Object().Value("code").Equal(http.StatusBadRequest)
}

func TestSignIn(t *testing.T) {
	e := expect(t)

	m := e.POST("/sign_in").WithJSON(testUser).Expect().Status(http.StatusOK).JSON().Object().ContainsKey("body")
	token = m.Value("body").String().Raw()

	e.POST("/sign_in").WithJSON(Stranger).Expect().Status(http.StatusOK).JSON().Object().Value("code").Equal(http.StatusBadRequest)
}

func TestPingAuth(t *testing.T) {
	e := expect(t)

	root := e.POST("/auth_ping").WithHeader("Authorization", "Bearer "+token).WithJSON(testUser).Expect().Status(http.StatusOK).JSON().Object()
	root.Value("code").Equal(http.StatusOK)
	root.Value("body").Equal(testUser.Name)

	dependencies.DBInstance().Delete(&testUser)
}