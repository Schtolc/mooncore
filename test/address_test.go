package test

import (
	"fmt"
	"github.com/Schtolc/mooncore/dao"
	"math/rand"
	"net/http"
	"testing"
)

func TestGetAddress(t *testing.T) {
	e := expect(t)

	lat := rand.Float64() + 55
	lon := rand.Float64() + 37

	address, err := dao.CreateAddress(lat, lon)

	if err != nil {
		t.Error("address cannot be created")
	}

	reqParams := fmt.Sprintf("id:\"%v\"", address.ID)
	respParams := "id, lat, lon, description"
	query := graphQLBody("{address(%s){%s}}", reqParams, respParams)

	resp := e.POST(graphqlURL).
		WithBytes(query).Expect().
		Status(http.StatusOK).JSON().Object().Value("data").Object().Value("address").Object()
	resp.Value("lat").Equal(lat)
	resp.Value("lon").Equal(lon)

	if err := dao.DeleteAddress(address.ID); err != nil {
		t.Error("address cannot be deleted")
	}
}
