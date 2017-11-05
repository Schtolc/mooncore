package handlers

import (
	"github.com/Schtolc/mooncore/models"
)

var (
	testUser2 = &models.UserAuth{
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

