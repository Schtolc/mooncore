package test

import (
	"fmt"
	"github.com/Schtolc/mooncore/dao"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
)

func TestEditMaster(t *testing.T) {
	e := expect(t)

	email := getEmail()
	password, passwordHash := getPassword()
	master, err := dao.CreateMaster(email, passwordHash)
	if err != nil {
		t.Error("cannot create master")
	}
	defer dao.DeleteMaster(master.ID)

	reqParamsAuth := fmt.Sprintf("email:\"%s\", password:\"%s\"", email, password)
	respParamsAuth := "token"
	queryAuth := graphQLBody("mutation {signIn(%s){%s}}", reqParamsAuth, respParamsAuth)

	token := e.POST(graphqlURL).
		WithBytes(queryAuth).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").
		Object().Value("signIn").Object().ContainsKey("token").Value("token").String().Raw()

	name := randString(10)
	photo := randString(10)
	latString := fmt.Sprintf("%f", rand.Float64())
	lonString := fmt.Sprintf("%f", rand.Float64())
	lat, _ := strconv.ParseFloat(latString, 64)
	lon, _ := strconv.ParseFloat(lonString, 64)
	reqParams := fmt.Sprintf("name:\"%s\", photo:\"%s\", lat:\"%s\", lon:\"%s\"", name, photo, latString, lonString)
	respParams := "id, name, avatar{id, path}, address{id, lat, lon}"
	query := graphQLBody("mutation {editMaster(%s){%s}}", reqParams, respParams)

	result := e.POST(graphqlURL).
		WithBytes(query).WithHeader("Authorization", "Bearer "+token).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").
		Object().Value("editMaster").Object()
	result.ContainsKey("name").Value("name").String().Equal(name)
	avatar := result.ContainsKey("avatar").Value("avatar").Object()
	avatar.ContainsKey("path").Value("path").String().Equal(photo)
	address := result.ContainsKey("address").Value("address").Object()
	address.ContainsKey("lat").Value("lat").Equal(lat)
	address.ContainsKey("lon").Value("lon").Equal(lon)
}
