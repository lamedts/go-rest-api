package proxy

import (
	"go-rest-api/types/datastore"
	"go-rest-api/utility/googleapi"
)

type Proxy struct {
	db     *datastore.DataStore
	mapapi *googlemapapi.GoogleMapAPI
}

func NewProxy(db datastore.DataStore, mapapi *googlemapapi.GoogleMapAPI) *Proxy {
	proxy := Proxy{
		db:     &db,
		mapapi: mapapi,
	}
	return &proxy
}
