package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/gavv/httpexpect"

	"github.com/Schtolc/mooncore/models"

	"math/rand"
	"net/http"
	"net/url"
	"testing"
)

var (
	conf      = dependencies.ConfigInstance()
	localhost = url.URL{Scheme: "http", Host: conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port}
	testUser  = &models.User{
		Email:    "newEmail",
		Password: "passPass",
		Role:     0,
	}
	testUserDetails = &models.UserDetails{
		UserID:    1,
		Name:      "aptrik",
		AddressID: 1,
		PhotoID:   1,
	}
)

type graphQLQuery struct {
	Query string `json:"query"`
}

func graphQLBody(query string, a ...interface{}) []byte {
	body, _ := json.Marshal(graphQLQuery{
		fmt.Sprintf(query, a...),
	})
	return body
}

func expect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  localhost.String(),
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: nil,
	})
}

// normal
// bad param - no such value
// bad param - invalid value
// without param

func TestCreateUser(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", password: \"%s\", role: %d", testUser.Email, testUser.Password, testUser.Role)
	respParams := "id, email, role"
	query := graphQLBody("mutation{createUser(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	body := resp.Value("body").Object().Value("createUser").Object()
	body.Value("email").Equal(testUser.Email)
	body.Value("role").Equal(testUser.Role)
}

func TestCreateUserBadParams(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", password: %s, role: %d", testUser.Email, testUser.Password, testUser.Role)
	respParams := "id, email, role"
	query := graphQLBody("mutation{createUser(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	errorMessage := "Argument \"password\" has invalid value passPass.\nExpected type \"String\", found passPass."
	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)
}

func TestCreateUserWithoutParams(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", password: \"%s\"", testUser.Email, testUser.Password)
	respParams := "id, email, role"
	query := graphQLBody("mutation{createUser(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	errorMessage := "Field \"createUser\" argument \"role\" of type \"Int!\" is required but not provided."
	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)
}

func TestCreateAddress(t *testing.T) {
	e := expect(t)

	address := &models.Address{
		Lat:         rand.Float64(),
		Lon:         rand.Float64(),
		Description: "description",
	}

	reqParams := fmt.Sprintf("lat:\"%.16v\", lon:\"%.16v\", description:\"%s\"", address.Lat, address.Lon, address.Description)
	respParams := "id, lat, lon, description"
	query := graphQLBody("mutation{createAddress(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	resp.Value("code").Number().Equal(http.StatusOK)

	body := resp.Value("body").Object().Value("createAddress")
	body.Object().Value("lat").Equal(address.Lat)
	body.Object().Value("lon").Equal(address.Lon)
	body.Object().Value("description").Equal(address.Description)
}

func TestCreateAddressBadParamLat(t *testing.T) {
	e := expect(t)

	address := &models.Address{
		Lat:         rand.Float64(),
		Lon:         rand.Float64(),
		Description: "description",
	}

	reqParams := fmt.Sprintf("lat:\"string\", lon:\"string\", description:\"%s\"", address.Description)
	respParams := "id, lat, lon, description"
	query := graphQLBody("mutation{createAddress(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal("InvalidParam: lat")
}

func TestCreateAddressBadParamDescription(t *testing.T) {
	e := expect(t)

	address := &models.Address{
		Lat:         rand.Float64(),
		Lon:         rand.Float64(),
		Description: "description",
	}

	reqParams := fmt.Sprintf("lat:\"%.16v\", lon:\"%.16v\", description:%.16v", address.Lat, address.Lon, address.Lon)
	respParams := "id, lat, lon, description"
	query := graphQLBody("mutation{createAddress(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	errorMessage := fmt.Sprintf("Argument \"description\" has invalid value 0.6868230728671094.\nExpected type \"String\", found %.16v.", address.Lon)

	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)
}

func TestCreateAddressWithoutParams(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("lat:\"string\", lon:\"string\"")
	respParams := "id, lat, lon, description"
	query := graphQLBody("mutation{createAddress(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	errorMessage := "Field \"createAddress\" argument \"description\" of type \"String!\" is required but not provided."
	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)

}

func TestCreatePhoto(t *testing.T) {
	e := expect(t)

	path := "random_path"
	tags := []int{1, 2}

	reqParams := fmt.Sprintf("path:\"%s\", tags:[%d,%d]", path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := graphQLBody("mutation{createPhoto(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")
	resp.Value("code").Number().Equal(http.StatusOK)

	body := resp.Value("body").Object().Value("createPhoto").Object()
	body.Value("path").Equal(path)
	body.Value("tags").Array().First().Object().Value("id").Equal(tags[0])
	body.Value("tags").Array().Last().Object().Value("id").Equal(tags[1])
}

func TestCreatePhotoWithNotExistingTags(t *testing.T) {
	e := expect(t)

	path := "random_path"
	tags := []int{1000, 2000}
	reqParams := fmt.Sprintf("path:\"%s\", tags:[%d,%d]", path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := graphQLBody("mutation{createPhoto(%s){%s}}", reqParams, respParams)
	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")
	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal("No such tags")
}

func TestCreatePhotoWithBadParamsTags(t *testing.T) {
	e := expect(t)

	path := "random_path"
	tags := []float64{
		rand.Float64(), rand.Float64(),
	}

	reqParams := fmt.Sprintf("path:\"%s\", tags:[%0.16v,%0.16v]", path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := graphQLBody("mutation{createPhoto(%s){%s}}", reqParams, respParams)
	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	errorMessage := fmt.Sprintf("Argument \"tags\" has invalid value [%0.16v, %0.16v].\nIn element #1: Expected type \"Int\", found %0.16v.\nIn element #1: Expected type \"Int\", found %0.16v.", tags[0], tags[1], tags[0], tags[1])

	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)
}

func TestCreatePhotoWithoutTags(t *testing.T) {
	e := expect(t)

	path := "random_path"

	reqParams := fmt.Sprintf("path:\"%s\"", path)
	respParams := "path, tags{id, name}"
	query := graphQLBody("mutation{createPhoto(%s){%s}}", reqParams, respParams)
	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	errorMessage := "Field \"createPhoto\" argument \"tags\" of type \"[Int]!\" is required but not provided."
	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)
}

func TestCreateSign(t *testing.T) {
	e := expect(t)

	signs := []int{1, 2}
	reqParams := fmt.Sprintf("email:\"%s\", signs:[%d, %d]", testUser.Email, signs[0], signs[1])
	respParams := "id, signs{id, name, photo{path}, description}"
	query := graphQLBody("mutation{addSigns(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	resp.Value("code").Number().Equal(http.StatusOK)

	body := resp.Value("body").Object().Value("addSigns").Object()
	body.Value("signs").Array().First().Object().Value("description").NotNull()
	body.Value("signs").Array().First().Object().Value("id").NotNull()
	body.Value("signs").Array().First().Object().Value("name").NotNull()
	body.Value("signs").Array().First().Object().Value("photo").Object().Value("path").NotNull()

	body.Value("signs").Array().Last().Object().Value("description").NotNull()
	body.Value("signs").Array().Last().Object().Value("id").NotNull()
	body.Value("signs").Array().Last().Object().Value("name").NotNull()
	body.Value("signs").Array().Last().Object().Value("photo").Object().Value("path").NotNull()

}

func TestCreateSignWithBadParamsSigns(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\"", testUser.Email)
	respParams := "id, signs{id, name}"
	query := graphQLBody("mutation{addSigns(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	errorMessage := "Field \"addSigns\" argument \"signs\" of type \"[Int]!\" is required but not provided."
	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)

}

func TestCreateSignWithoutSigns(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\"", testUser.Email)
	respParams := "id, signs{id, name}"
	query := graphQLBody("mutation{addSigns(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	errorMessage := "Field \"addSigns\" argument \"signs\" of type \"[Int]!\" is required but not provided."
	resp.Value("code").Number().Equal(http.StatusNotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)

}

func TestCreateUserProfile(t *testing.T) {
	e := expect(t)
	photos := []int{1, 2}
	signs := []int{1, 2}
	reqParams := fmt.Sprintf("email:\"%s\", name:\"%s\", address_id: %d, avatar_id: %d, photos: [%d, %d], signs: [%d, %d]",
		testUser.Email, testUserDetails.Name, testUserDetails.AddressID, testUserDetails.PhotoID, photos[0], photos[1], signs[0], signs[1])
	respParams := fmt.Sprintf("id, name, address{lat,lon,id}, avatar{id, path}, photos{id, path}, signs{id, name}")

	query := graphQLBody("mutation{createUserProfile(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	resp.Value("code").Number().Equal(http.StatusOK)
	body := resp.Value("body").Object().Value("createUserProfile").Object()

	body.Value("name").Equal(testUserDetails.Name)
	body.Value("address").Object().Value("id").Equal(testUserDetails.AddressID)
	body.Value("address").Object().Value("lat").NotNull()
	body.Value("address").Object().Value("lon").NotNull()

	body.Value("avatar").Object().Value("id").Equal(testUserDetails.PhotoID)
	body.Value("avatar").Object().Value("path").NotNull()

	body.Value("photos").Array().First().Object().Value("id").Equal(1)
	body.Value("photos").Array().Last().Object().Value("id").Equal(2)
	body.Value("photos").Array().First().Object().Value("path").NotNull()
	body.Value("photos").Array().Last().Object().Value("path").NotNull()

	body.Value("signs").Array().First().Object().Value("id").Equal(1)
	body.Value("signs").Array().Last().Object().Value("id").Equal(2)
	body.Value("signs").Array().First().Object().Value("name").NotNull()
	body.Value("signs").Array().Last().Object().Value("name").NotNull()
}

//func TestGetAddress(t *testing.T) {
//	e := expect(t)
//
//	address := createAddress(t)
//
//	query := fmt.Sprintf("{address(id:%d){lat, lon}}", address.ID)
//
//	root := e.GET("/graphql").
//		WithQuery("query", query).Expect().
//		Status(http.StatusOK).JSON()
//
//	root.Object().Value("code").Number().Equal(OK)
//
//	obj := root.Object().Value("body").
//		Object().Value("address")
//
//	lat := obj.Object().Value("lat").Number().Raw()
//	lon := obj.Object().Value("lon").Number().Raw()
//
//	assert.Equal(t, int(address.Lat), int(lat), "created lat doesn't equal to returned")
//	assert.Equal(t, int(address.Lon), int(lon), "created lon doesn't equal to returned")
//
//	db.Delete(&address)
//
//	root = e.GET("/graphql").
//		WithQuery("query", query).Expect().
//		Status(http.StatusOK).JSON()
//
//	root.Object().Value("code").Number().Equal(NotFound)
//	root.Object().Value("body").NotNull()
//}
//

//
//func TestGetPhoto(t *testing.T) {
//	e := expect(t)
//
//	photo := createPhoto(t)
//
//	query := fmt.Sprintf("{photo(id:%d){path}}", photo.ID)
//
//	root := e.GET("/graphql").
//		WithQuery("query", query).Expect().
//		Status(http.StatusOK).JSON()
//
//	root.Object().Value("code").Number().Equal(OK)
//
//	root = root.Object().Value("body").
//		Object().Value("photo")
//
//	path := root.Object().Value("path").Raw()
//
//	assert.Equal(t, photo.Path, path, "created path doesn't equal to returned")
//
//	db.Delete(&photo)
//
//	root = e.GET("/graphql").
//		WithQuery("query", query).Expect().
//		Status(http.StatusOK).JSON()
//
//	root.Object().Value("code").Number().Equal(NotFound)
//	root.Object().Value("body").NotNull()
//}
//
//func TestCreateUser(t *testing.T) {
//	e := expect(t)
//
//	address := createAddress(t)
//	photo := createPhoto(t)
//
//	name := randString()
//	password := randString()
//	email := randString()
//	addressID := address.ID
//	photoID := photo.ID
//
//	query := fmt.Sprintf("mutation{createUser(name:\"%s\", "+
//		"password:\"%s\", "+
//		"email:\"%s\", "+
//		"address_id:%d, "+
//		"photo_id:%d){id}}", name, password, email, addressID, photoID)
//
//	root := e.GET("/graphql").
//		WithQuery("query", query).Expect().
//		Status(http.StatusOK).JSON()
//
//	root.Object().Value("code").Number().Equal(OK)
//
//	id := root.Object().Value("body").
//		Object().Value("createUser").
//		Object().Value("id").Number().Raw()
//
//	user := models.User{}
//
//	assert.Nil(t, db.First(&user, int(id)).Error, "user was not created")
//
//	assert.Equal(t, user.Name, name, "created name doesn't equal to returned")
//	assert.Equal(t, user.Password, password, "created password doesn't equal to returned")
//	assert.Equal(t, user.Email, email, "created email doesn't equal to returned")
//	assert.Equal(t, user.AddressID, addressID, "created addressID doesn't equal to returned")
//	assert.Equal(t, user.PhotoID, photoID, "created photoID doesn't equal to returned")
//
//	db.Delete(&user)
//	db.Delete(&photo)
//	db.Delete(&address)
//}
//
//func TestGetUser(t *testing.T) {
//	e := expect(t)
//
//	address := createAddress(t)
//	photo := createPhoto(t)
//
//	user := models.User{
//		Name:      randString(),
//		Password:  randString(),
//		Email:     randString(),
//		AddressID: address.ID,
//		PhotoID:   photo.ID,
//	}
//
//	assert.Nil(t, db.Create(&user).Error, "user was not created")
//
//	query := fmt.Sprintf("{user(id:%d){name, email, address{id}, photo{id}}}", user.ID)
//
//	root := e.GET("/graphql").
//		WithQuery("query", query).Expect().
//		Status(http.StatusOK).JSON()
//
//	root.Object().Value("code").Number().Equal(OK)
//
//	root = root.Object().Value("body").
//		Object().Value("user")
//
//	root.Object().Value("name").Equal(user.Name)
//	root.Object().Value("email").Equal(user.Email)
//	root.Object().Value("address").Object().Value("id").Equal(user.AddressID)
//	root.Object().Value("photo").Object().Value("id").Equal(user.PhotoID)
//
//	db.Delete(&user)
//	db.Delete(&photo)
//	db.Delete(&address)
//
//	root = e.GET("/graphql").
//		WithQuery("query", query).Expect().
//		Status(http.StatusOK).JSON()
//
//	root.Object().Value("code").Number().Equal(NotFound)
//	root.Object().Value("body").NotNull()
//}
