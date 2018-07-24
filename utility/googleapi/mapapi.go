package googlemapapi

import (
	"errors"
	"go-rest-api/logger"
	"time"

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

// mapAPI.Dirctions([2]float32{2, 40}, [2]float32{2, 40})
func (googleMapAPI *GoogleMapAPI) Dirctions(originCoord string, destCoord string) (*time.Duration, *int, error) {

	r := &maps.DirectionsRequest{
		Origin:      originCoord,
		Destination: destCoord,
	}
	routes, _, err := googleMapAPI.client.Directions(context.Background(), r)
	if err != nil {
		googleMapAPILogger.Errorf("fatal error: %s", err)
		return nil, nil, err
	}

	// A route with no waypoints will contain exactly one leg within the legs array.
	// Route represents a single route between an origin and a destination.
	//
	// pretty.Println(routes[0].Legs[0].Distance)
	// pretty.Println(routes[0].Legs[0].Duration)
	if len(routes) == 1 && len(routes[0].Legs) == 1 {
		return &routes[0].Legs[0].Duration, &routes[0].Legs[0].Distance.Meters, nil
	}
	return nil, nil, errors.New("Unexpected response from google map api")
}
