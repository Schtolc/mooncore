package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Schtolc/mooncore/dependencies"
	"github.com/Schtolc/mooncore/models"
	"github.com/gavv/httpexpect"
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
		Role:     1,
	}
	testUserDetails = &models.UserDetails{
		UserID:    1,
		Name:      "user_name",
		AddressID: 1,
		PhotoID:   1,
	}
)

type graphQLQuery struct {
	Query string `json:"query"`
}

func randString() string {
	var result string
	l := rand.Int() % 50
	for i := 0; i < l; i++ {
		result += string(rand.Int()%('z'-'a') + 'a')
	}
	return result
}

func graphQLBody(query string, a ...interface{}) []byte {
	body, _ := json.Marshal(graphQLQuery{
		fmt.Sprintf(query, a...),
	})
	return body
}

func expect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  localhost.String() + conf.Server.APIPrefix,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: nil,
	})
}

// [[ TEST CREATE USER ]]

func TestCreateUser(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", password: \"%s\", role: %d", testUser.Email, testUser.Password, testUser.Role)
	respParams := "id, email, role"
	query := graphQLBody("mutation{createUser(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("createUser").Object()
	resp.Value("email").Equal(testUser.Email)
	resp.Value("role").Equal(testUser.Role)
}

func TestCreateSecondUser(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", password: \"%s\", role: %d", testUser.Email, testUser.Password, testUser.Role)
	respParams := "id, email, role"
	query := graphQLBody("mutation{createUser(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("createUser").Object()

	resp.Value("email").Equal(testUser.Email)
	resp.Value("role").Equal(testUser.Role)
}

func TestCreateUserBadParams(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", password: %s, role: %d", testUser.Email, testUser.Password, testUser.Role)
	respParams := "id, email, role"
	query := graphQLBody("mutation{createUser(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusNotFound).JSON().Object().Value("data").Array()
	errorMessage := "Argument \"password\" has invalid value passPass.\nExpected type \"String\", found passPass."
	resp.First().Object().Value("message").Equal(errorMessage)
}

func TestCreateUserWithoutParams(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", password: \"%s\"", testUser.Email, testUser.Password)
	respParams := "id, email, role"
	query := graphQLBody("mutation{createUser(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusNotFound).JSON().Object().Value("data").Array()

	errorMessage := "Field \"createUser\" argument \"role\" of type \"Int!\" is required but not provided."
	resp.First().Object().Value("message").Equal(errorMessage)
}

// [[ QUERY USER ]]

func TestGetUser(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\"", testUser.Email)
	respParams := "name, address{lat, lon}, avatar{path}"
	query := graphQLBody("{getUser(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("getUser").Object()
	resp.Value("address").Object().Value("lat").Equal(0)
	resp.Value("address").Object().Value("lon").Equal(0)
	resp.Value("avatar").Object().Value("path").Equal("default")
}

func TestListUsers(t *testing.T) {
	e := expect(t)

	respParams := "id, name, address{lat, lon}, avatar{path}"
	query := graphQLBody("{listUsers{%s}}", respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("listUsers").Array()

	resp.Element(1).Object().Value("name").Equal("")
	resp.Element(1).Object().Value("address").Object().Value("lat").Equal(0)
	resp.Element(1).Object().Value("address").Object().Value("lon").Equal(0)
	resp.Element(1).Object().Value("avatar").Object().Value("path").Equal("default")
}

// [[ TEST USER DETAILS EDIT]]

func TestEditUserProfile(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", name:\"%s\", address_id: %d, avatar_id: %d",
		testUser.Email, testUserDetails.Name, testUserDetails.AddressID, testUserDetails.PhotoID)
	respParams := fmt.Sprintf("id, name, address{lat,lon,id}, avatar{id, path}, photos{id, path}")

	query := graphQLBody("mutation{editUserProfile(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("editUserProfile").Object()

	resp.Value("name").Equal(testUserDetails.Name)
	resp.Value("address").Object().Value("id").Equal(testUserDetails.AddressID)
	resp.Value("address").Object().Value("lat").NotNull()
	resp.Value("address").Object().Value("lon").NotNull()

	resp.Value("avatar").Object().Value("id").Equal(testUserDetails.PhotoID)
	resp.Value("avatar").Object().Value("path").NotNull()
}

func TestEditUserProfileNotAllParams(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\", avatar_id: %d",
		testUser.Email, testUserDetails.PhotoID)
	respParams := fmt.Sprintf("id, name, address{lat,lon,id}, avatar{id, path}, photos{id, path}")

	query := graphQLBody("mutation{editUserProfile(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("editUserProfile").Object()

	resp.Value("name").Equal(testUserDetails.Name)
	resp.Value("address").Object().Value("id").Equal(testUserDetails.AddressID)
	resp.Value("address").Object().Value("lat").NotNull()
	resp.Value("address").Object().Value("lon").NotNull()

	resp.Value("avatar").Object().Value("id").Equal(testUserDetails.PhotoID)
	resp.Value("avatar").Object().Value("path").NotNull()
}
func TestCreateAddress(t *testing.T) {
	e := expect(t)

	address := &models.Address{
		Lat:         rand.Float64(),
		Lon:         rand.Float64(),
		Description: "description",
	}

	reqParams := fmt.Sprintf("lat:\"%.20v\", lon:\"%.20v\", description:\"%s\"", address.Lat, address.Lon, address.Description)
	respParams := "id, lat, lon, description"
	query := graphQLBody("mutation{createAddress(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("createAddress").Object()

	resp.Value("lat").Equal(address.Lat)
	resp.Value("lon").Equal(address.Lon)
	resp.Value("description").Equal(address.Description)
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
		Status(http.StatusNotFound).JSON().Object().Value("data").Array()

	errorMessage := "strconv.ParseFloat: parsing \"string\": invalid syntax"
	resp.First().Object().Value("message").Equal(errorMessage)
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
		Status(http.StatusNotFound).JSON().Object().Value("data").Array()
	errorMessage := fmt.Sprintf("Argument \"description\" has invalid value %.16v.\nExpected type \"String\", found %.16v.", address.Lon, address.Lon)
	resp.First().Object().Value("message").Equal(errorMessage)
}

func TestCreateAddressWithoutParams(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("lat:\"string\", lon:\"string\"")
	respParams := "id, lat, lon, description"
	query := graphQLBody("mutation{createAddress(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusNotFound).JSON().Object().Value("data").Array()

	errorMessage := "Field \"createAddress\" argument \"description\" of type \"String!\" is required but not provided."
	resp.First().Object().Value("message").Equal(errorMessage)
}

func TestCreateAvatar(t *testing.T) {
	e := expect(t)

	path := "random_path"
	tags := []int{1, 2}

	reqParams := fmt.Sprintf("path:\"%s\", tags:[%d,%d]", path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := graphQLBody("mutation{createAvatar(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("createAvatar").Object()

	resp.Value("path").Equal(path)
	resp.Value("tags").Array().First().Object().Value("id").Equal(tags[0])
	resp.Value("tags").Array().Last().Object().Value("id").Equal(tags[1])
}

func TestCreateAvatarWithNotExistingTags(t *testing.T) {
	e := expect(t)

	path := "random_path"
	tags := []int{1000, 2000}
	reqParams := fmt.Sprintf("path:\"%s\", tags:[%d,%d]", path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := graphQLBody("mutation{createAvatar(%s){%s}}", reqParams, respParams)
	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusNotFound).JSON().Object().Value("data").Array()

	resp.First().Object().Value("message").Equal("No such tags")
}

func TestCreateAvatarWithBadParamsTags(t *testing.T) {
	e := expect(t)

	path := "random_path"
	tags := []float64{
		rand.Float64(), rand.Float64(),
	}

	reqParams := fmt.Sprintf("path:\"%s\", tags:[%0.16v,%0.16v]", path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := graphQLBody("mutation{createAvatar(%s){%s}}", reqParams, respParams)
	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusNotFound).JSON().Object().Value("data").Array()

	errorMessage := fmt.Sprintf("Argument \"tags\" has invalid value [%0.16v, %0.16v].\nIn element #1: Expected type \"Int\", found %0.16v.\nIn element #1: Expected type \"Int\", found %0.16v.", tags[0], tags[1], tags[0], tags[1])

	resp.First().Object().Value("message").Equal(errorMessage)
}

func TestCreateAvatarWithoutTags(t *testing.T) {
	e := expect(t)

	path := "random_path"

	reqParams := fmt.Sprintf("path:\"%s\"", path)
	respParams := "path"
	query := graphQLBody("mutation{createAvatar(%s){%s}}", reqParams, respParams)
	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("createAvatar").Object()

	resp.Value("path").Equal(path)
}

// [[ TEST PHOTOS CREATE]]

func TestAddPhoto(t *testing.T) {
	e := expect(t)

	path := "random_path123"
	tags := []int{1, 2}

	reqParams := fmt.Sprintf("email:\"%s\", path:\"%s\", tags:[%d,%d]", testUser.Email, path, tags[0], tags[1])
	respParams := "path, tags{id, name}"
	query := graphQLBody("mutation{addPhoto(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("addPhoto").Object()
	resp.Value("path").Equal(path)
	resp.Value("tags").Array().First().Object().Value("id").Equal(tags[0])
	resp.Value("tags").Array().Last().Object().Value("id").Equal(tags[1])
}

// [[ TEST PHOTOS GET]]

func TestGetPhoto(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("id:2")
	respParams := "path, tags{id, name}"
	query := graphQLBody("{getPhoto(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("getPhoto").Object()
	resp.Value("path").NotNull()
	resp.Value("tags").Array().First().Object().Value("id").NotNull()
	resp.Value("tags").Array().Last().Object().Value("id").NotNull()
}

func TestGetUserPhotos(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\"", testUser.Email)
	respParams := "path, tags{id, name}"
	query := graphQLBody("{getUserPhotos(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("getUserPhotos").Array()

	resp.First().Object().Value("path").NotNull()
	resp.First().Object().Value("tags").NotNull()
}

// [[ TEST ADD SIGN ]]

func TestCreateSign(t *testing.T) {
	e := expect(t)

	signs := []int{1, 2}
	reqParams := fmt.Sprintf("email:\"%s\", signs:[%d, %d]", testUser.Email, signs[0], signs[1])
	respParams := "id, signs{id, name, icon, description}"
	query := graphQLBody("mutation{addSigns(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("addSigns").Object()

	firstElem := resp.Value("signs").Array().First().Object()
	firstElem.Value("description").NotNull()
	firstElem.Value("id").NotNull()
	firstElem.Value("name").NotNull()
	firstElem.Value("icon").NotNull()

	lastElem := resp.Value("signs").Array().Last().Object()
	lastElem.Value("description").NotNull()
	lastElem.Value("id").NotNull()
	lastElem.Value("name").NotNull()
	lastElem.Value("icon").NotNull()
}

func TestCreateSignWithBadParamsSigns(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\"", testUser.Email)
	respParams := "id, signs{id, name}"
	query := graphQLBody("mutation{addSigns(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusNotFound).JSON().Object().Value("data").Array()

	errorMessage := "Field \"addSigns\" argument \"signs\" of type \"[Int]!\" is required but not provided."
	resp.First().Object().Value("message").Equal(errorMessage)

}

func TestCreateSignWithoutSigns(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\"", testUser.Email)
	respParams := "id, signs{id, name}"
	query := graphQLBody("mutation{addSigns(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusNotFound).JSON().Object().Value("data").Array()

	errorMessage := "Field \"addSigns\" argument \"signs\" of type \"[Int]!\" is required but not provided."
	resp.First().Object().Value("message").Equal(errorMessage)
}

// [[ TEST GET SIGN ]]

func TestGetSign(t *testing.T) {
	e := expect(t)

	reqParams := fmt.Sprintf("email:\"%s\"", testUser.Email)
	respParams := "id, name, icon, description"
	query := graphQLBody("{getSigns(%s){%s}}", reqParams, respParams)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object()

	resp.Keys().ContainsOnly("code", "body")

	resp.Value("code").Number().Equal(http.StatusOK)

	body := resp.Value("body").Object().Value("getSigns").Array()
	body.First().Object().Value("name").Equal("accuracy")
	body.First().Object().Value("description").Equal("means accuracy")

	body.First().Object().Value("icon").Equal("default")
}

//// [[ GET FEED ]]

func TestGetFeed(t *testing.T) {
	e := expect(t)
	query := graphQLBody("{feed(limit:3){id, name, avatar{id, path, tags{id, name}}, photos{id, path, tags{name}}, signs{id, name, description, icon}}}")
	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("feed")
	photo := resp.Array().First().Object().Value("photos").Array().First()
	photo.Object().Value("id").NotNull()
	photo.Object().Value("path").NotNull()
	photo.Object().Value("tags").Array()

	sign := resp.Array().First().Object().Value("signs").Array().First()
	sign.Object().Value("description").NotNull()
	sign.Object().Value("name").NotNull()
	sign.Object().Value("id").NotNull()
	sign.Object().Value("icon").NotNull()
	//sign.Object().Value("photo").Object().Value("tags").Array()

	pho := resp.Array().First().Object().Value("avatar")
	pho.Object().Value("id").NotNull()
	pho.Object().Value("path").NotNull()
	//pho.Object().Value("tags").Array().First()
}

//// [[ GET ADDRESSES IN GIVEN AREA ]]

func TestGetAddressInGivenArea(t *testing.T) {
	e := expect(t)

	db := dependencies.DBInstance()

	boundary1 := models.Address{
		Lat: 5,
		Lon: 10,
	}
	boundary2 := models.Address{
		Lat: 10,
		Lon: 5,
	}

	halfSide := (boundary2.Lat - boundary1.Lat) / 2
	latInArea := (boundary1.Lat + boundary2.Lat) / 2
	lonInArea := (boundary1.Lon + boundary2.Lon) / 2

	objectInArea := models.Address{
		Lat: latInArea,
		Lon: lonInArea,
	}

	objectOutOfArea1 := models.Address{
		Lat: boundary1.Lat - halfSide,
		Lon: lonInArea,
	}
	objectOutOfArea2 := models.Address{
		Lat: boundary2.Lat + halfSide,
		Lon: lonInArea,
	}
	objectOutOfArea3 := models.Address{
		Lat: latInArea,
		Lon: boundary1.Lon + halfSide,
	}
	objectOutOfArea4 := models.Address{
		Lat: latInArea,
		Lon: boundary2.Lon - halfSide,
	}

	db.Create(&objectInArea)
	db.Create(&objectOutOfArea1)
	db.Create(&objectOutOfArea2)
	db.Create(&objectOutOfArea3)
	db.Create(&objectOutOfArea4)

	query := graphQLBody(
		"{addressListInArea(lat1:\"%f\", lon1:\"%f\", lat2:\"%f\", lon2:\"%f\"){lat, lon, id}}",
		boundary1.Lat,
		boundary1.Lon,
		boundary2.Lat,
		boundary2.Lon,
	)

	resp := e.POST("/graphql").
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("addressListInArea").Array()
	db.Delete(&objectInArea)
	db.Delete(&objectOutOfArea1)
	db.Delete(&objectOutOfArea2)
	db.Delete(&objectOutOfArea3)
	db.Delete(&objectOutOfArea4)

	resp.Length().Equal(1)
	resp.Element(0).Object().Value("id").Equal(objectInArea.ID)
}
