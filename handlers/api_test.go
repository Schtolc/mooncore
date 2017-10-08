package handlers

import (
	"fmt"
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/gavv/httpexpect"

	"github.com/Schtolc/mooncore/models"

	"net/url"
	"testing"
	"math/rand"
	"net/http"

	"github.com/sirupsen/logrus"
)

var (
	conf      = dependencies.ConfigInstance()
	localhost = url.URL{Scheme: "http", Host: conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port}
	//db        = dependencies.DBInstance()
)

func expect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  localhost.String(),
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: nil,
	})
}

//func randString() string {
//	var result string
//	l := rand.Int() % 50
//	for i := 0; i < l; i++ {
//		result += string(rand.Int()%('z'-'a') + 'a')
//	}
//	return result
//}
//
//func createAddress(t *testing.T) *models.Address {
//	address := &models.Address{
//		Lat: rand.Float64(),
//		Lon: rand.Float64(),
//	}
//	assert.Nil(t, db.Create(&address).Error, "address was not created")
//	return address
//}
//
//func createPhoto(t *testing.T) *models.Photo {
//	photo := &models.Photo{
//		Path: randString(),
//	}
//	assert.Nil(t, db.Create(&photo).Error, "photo was not created")
//	return photo
//}

func TestCreateAddress (t *testing.T) {
	e := expect(t)

	address := &models.Address{
		Lat: rand.Float64(),
		Lon: rand.Float64(),
		Description: "description",
	}

	reqParams  := fmt.Sprintf("lat:\"%.16v\", lon:\"%.16v\", description:\"%s\"",address.Lat, address.Lon, address.Description)
	respParams := "id, lat, lon, description"
	query := fmt.Sprintf("mutation{createAddress(%s){%s}}",reqParams, respParams)

	resp := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	resp.Value("code").Number().Equal(OK)

	body := resp.Value("body").Object().Value("createAddress")
	body.Object().Value("lat").Equal(address.Lat)
	body.Object().Value("lon").Equal(address.Lon)
	body.Object().Value("description").Equal(address.Description)
}

func TestCreateAddressBadParamLat (t *testing.T) {
	e := expect(t)

	address := &models.Address{
		Lat: rand.Float64(),
		Lon: rand.Float64(),
		Description: "description",
	}

	reqParams  := fmt.Sprintf("lat:\"string\", lon:\"string\", description:\"%s\"", address.Description)
	respParams := "id, lat, lon, description"
	query := fmt.Sprintf("mutation{createAddress(%s){%s}}",reqParams, respParams)

	resp := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Value("code").Number().Equal(NotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal("InvalidParam: lat")
}

func TestCreateAddressBadParamDescription (t *testing.T) {
	e := expect(t)

	address := &models.Address{
		Lat: rand.Float64(),
		Lon: rand.Float64(),
		Description: "description",
	}

	reqParams  := fmt.Sprintf("lat:\"%.16v\", lon:\"%.16v\", description:%.16v",address.Lat, address.Lon, address.Lon)
	respParams := "id, lat, lon, description"
	query := fmt.Sprintf("mutation{createAddress(%s){%s}}",reqParams, respParams)

	resp := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().Object()

	errorMessage := fmt.Sprintf("Argument \"description\" has invalid value 0.6868230728671094.\nExpected type \"String\", found %.16v.", address.Lon)

	resp.Value("code").Number().Equal(NotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)
}

func TestCreateAddressWithoutParams (t *testing.T) {
	e := expect(t)

	reqParams  := fmt.Sprintf("lat:\"string\", lon:\"string\"")
	respParams := "id, lat, lon, description"
	query := fmt.Sprintf("mutation{createAddress(%s){%s}}",reqParams, respParams)

	resp := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Value("code").Number().Equal(NotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal("Field \"createAddress\" argument \"description\" of type \"String!\" is required but not provided.")

}

func TestCreatePhoto(t *testing.T) {
	e := expect(t)

	path := "random_path"
	tags := []int{1,2}
	reqParams  := fmt.Sprintf("path:\"%s\", tags:[%d,%d]",path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := fmt.Sprintf("mutation{createPhoto(%s){%s}}",reqParams, respParams)
	resp := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")
	resp.Value("code").Number().Equal(OK)

	body := resp.Value("body").Object().Value("createPhoto").Object()
	body.Value("path").Equal(path)
	body.Value("tags").Array().First().Object().Value("id").Equal(tags[0])
	body.Value("tags").Array().Last().Object().Value("id").Equal(tags[1])
}

func TestCreatePhotoWithNotExistingTags(t *testing.T) {
	e := expect(t)

	path := "random_path"
	tags := []int{1000,2000}
	reqParams  := fmt.Sprintf("path:\"%s\", tags:[%d,%d]",path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := fmt.Sprintf("mutation{createPhoto(%s){%s}}",reqParams, respParams)
	resp := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")
	resp.Value("code").Number().Equal(NotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal("No such tags")

}
func TestCreatePhotoWithBadTags(t *testing.T) {
	e := expect(t)

	path := "random_path"
	tags := []float64{
		rand.Float64(), rand.Float64(),
	}

	reqParams  := fmt.Sprintf("path:\"%s\", tags:[%0.16v,%0.16v]",path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := fmt.Sprintf("mutation{createPhoto(%s){%s}}",reqParams, respParams)
	resp := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	errorMessage := fmt.Sprintf("Argument \"tags\" has invalid value [%0.16v, %0.16v].\nIn element #1: Expected type \"Int\", found %0.16v.\nIn element #1: Expected type \"Int\", found %0.16v.", tags[0], tags[1], tags[0], tags[1])

	resp.Value("code").Number().Equal(NotFound)
	resp.Value("body").Array().First().Object().Value("message").Equal(errorMessage)

}

func TestCreateSign(t *testing.T) {
	e := expect(t)

	signs := []int{1,2}
	id := 1
	reqParams  := fmt.Sprintf("id:%d, signs:[%d, %d]",id, signs[0],signs[1])
	respParams := "id, signs{id, name, path, description}"
	query := fmt.Sprintf("mutation{createSign(%s){%s}}",reqParams, respParams)
	logrus.Warn(query)
	resp := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	resp.Value("code").Number().Equal(OK)
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
