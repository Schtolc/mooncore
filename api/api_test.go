package api

import (
	"github.com/Schtolc/mooncore/utils"
	"net/url"
	"net/http"
	"testing"
	"github.com/gavv/httpexpect"
	"github.com/Schtolc/mooncore/database"
	"github.com/Schtolc/mooncore/models"
	"math/rand"
	"fmt"
	"math"
)

var (
	conf      = utils.GetConfig()
	localhost = url.URL{Scheme: "http", Host: conf.Server.Hostbase.Host + ":" + conf.Server.Hostbase.Port}
	db        = database.InitDatabase(conf)
)

func randString() string {
	var result string
	l := rand.Int() % 50
	for i := 0; i < l; i++ {
		result += string(rand.Int()%('z'-'a') + 'a')
	}
	return result
}

func TestCreateAddress(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	lat := rand.Float64()
	lon := rand.Float64()

	query := fmt.Sprintf("mutation{createAddress(lat:%f,lon:%f){id}}", lat, lon)

	id := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().
		Object().Value("data").
		Object().Value("createAddress").
		Object().Value("id").Number().Raw()

	address := models.Address{}
	db.First(&address, int(id))

	if math.Abs(address.Lat-lat) > 1 || math.Abs(address.Lon-lon) > 1 {
		t.Fail()
	}

	db.Delete(&address)
}

func TestGetAddress(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	address := &models.Address{
		Lat: rand.Float64(),
		Lon: rand.Float64(),
	}

	db.Create(&address)

	query := fmt.Sprintf("{address(id:%d){lat, lon}}", address.ID)

	obj := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().
		Object().Value("data").
		Object().Value("address")

	lat := obj.Object().Value("lat").Raw()
	lon := obj.Object().Value("lon").Raw()

	if lat != address.Lat || lon != address.Lon {
		t.Fail()
	}

	db.Delete(&address)
}

func TestCreatePhoto(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	path := randString()

	query := fmt.Sprintf("mutation{createPhoto(path:\"%s\"){id}}", path)

	id := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().
		Object().Value("data").
		Object().Value("createPhoto").
		Object().Value("id").Number().Raw()

	photo := models.Photo{}
	db.First(&photo, int(id))

	if photo.Path != path {
		t.Fail()
	}

	db.Delete(&photo)
}

func TestGetPhoto(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	photo := &models.Photo{
		Path: randString(),
	}

	db.Create(&photo)

	query := fmt.Sprintf("{photo(id:%d){path}}", photo.ID)

	obj := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().
		Object().Value("data").
		Object().Value("photo")

	path := obj.Object().Value("path").Raw()

	if path != photo.Path {
		t.Fail()
	}

	db.Delete(&photo)
}

func TestCreateUser(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	address := &models.Address{
		Lat: rand.Float64(),
		Lon: rand.Float64(),
	}

	db.Create(&address) // create test address

	path := randString()

	photo := &models.Photo{
		Path: path,
	}

	db.Create(&photo) // create test photo

	name := randString()
	password := randString()
	email := randString()
	addressId := address.ID
	photoId := photo.ID

	query := fmt.Sprintf("mutation{createUser(name:\"%s\", password:\"%s\", email:\"%s\", address_id:%d, photo_id:%d){id}}", name, password, email, addressId, photoId)

	id := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().
		Object().Value("data").
		Object().Value("createUser").
		Object().Value("id").Number().Raw()

	user := models.User{}
	db.First(&user, int(id))

	if name != user.Name || password != user.Password || email != user.Email || addressId != user.AddressID || photoId != user.PhotoID {
		t.Fail()
	}

	db.Delete(&user)
	db.Delete(&photo)
	db.Delete(&address)
}

func TestGetUser(t *testing.T) {
	e := httpexpect.New(t, localhost.String())

	address := &models.Address{
		Lat: rand.Float64(),
		Lon: rand.Float64(),
	}

	db.Create(&address) // create test address

	path := randString()

	photo := &models.Photo{
		Path: path,
	}

	db.Create(&photo) // create test photo

	user := models.User{
		Name:      randString(),
		Password:  randString(),
		Email:     randString(),
		AddressID: address.ID,
		PhotoID:   photo.ID,
	}

	db.Create(&user) // create test user

	query := fmt.Sprintf("{user(id:%d){name, email, address{id}, photo{id}}}", user.ID)

	userRoot := e.GET("/graphql").
		WithQuery("query", query).Expect().
		Status(http.StatusOK).JSON().
		Object().Value("data").
		Object().Value("user")

	userRoot.Object().Value("name").Equal(user.Name)
	userRoot.Object().Value("email").Equal(user.Email)
	userRoot.Object().Value("address").Object().Value("id").Equal(user.AddressID)
	userRoot.Object().Value("photo").Object().Value("id").Equal(user.PhotoID)

	db.Delete(&user)
	db.Delete(&photo)
	db.Delete(&address)
}
