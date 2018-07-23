package googlemapapi

import (
	"go-rest-api/logger"

	"github.com/kr/pretty"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var googleMapAPILogger = logger.GetLogger("googlemapapi")

type GoogleMapAPI struct {
	key    string
	client *maps.Client
}

func NewGoogleMapAPI(key string) *GoogleMapAPI {
	googleMapAPI := GoogleMapAPI{key: key}
	if client, err := maps.NewClient(maps.WithAPIKey(key)); err != nil {
		googleMapAPILogger.Errorf("fatal error: %s", err)
	} else {
		googleMapAPILogger.Infof("Google api connected")
		googleMapAPILogger.Debugf("%+v", client)
		googleMapAPI.client = client
	}
	return &googleMapAPI
}

func (googleMapAPI *GoogleMapAPI) Dirctions(orig [2]float32, dest [2]float32) {

	r := &maps.DirectionsRequest{
		Origin:      "22.375036,114.194551",
		Destination: "22.387436,114.208592",
	}
	route, _, err := googleMapAPI.client.Directions(context.Background(), r)
	if err != nil {
		googleMapAPILogger.Errorf("fatal error: %s", err)
	}
	googleMapAPILogger.Infof("fatal error: %v", route)
	pretty.Println(route)
}
