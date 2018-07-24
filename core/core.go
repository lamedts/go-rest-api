package core

import (
	"go-rest-api/core/proxy"
	"go-rest-api/datastore"
	"go-rest-api/logger"
	"go-rest-api/server/api"
	"go-rest-api/types"
	typesDatastore "go-rest-api/types/datastore"
	"go-rest-api/utility/googleapi"
)

var coreLogger = logger.GetLogger("core")

type Yay struct {
	config  types.Config
	servers []types.Server
	proxy   *proxy.Proxy
}

var yayInstance *Yay

func GetYay(configuration *types.Config) *Yay {
	apiEnabled := false
	graphqlEnabled := false
	if configuration.ServerSettings.APIServerConfig != nil {
		apiEnabled = true
	}
	if configuration.ServerSettings.GraphqlServerConfig != nil {
		graphqlEnabled = true
	}
	coreLogger.Infof(".............    config    .............")
	coreLogger.Infof("db: %v", configuration.DatastoreSettings.Type)
	coreLogger.Infof("api enabled: %v", apiEnabled)
	coreLogger.Infof("graphql enabled: %v", graphqlEnabled)
	coreLogger.Infof("........................................")
	if yayInstance != nil {
		return yayInstance
	}

	coreLogger.Infoln("**      Database Connection      **")
	var dataStore typesDatastore.DataStore
	switch configuration.DatastoreSettings.Type {
	// case DatastoreTypeMongodb:
	// 	coreLogger.Infof("loading mongodb")
	// 	dataStore = datastore.NewMongoDB(configuration.DatastoreSettings)
	default:
		// TODO: verify config
		dataStore = datastore.NewMysqlDB(configuration.DatastoreSettings)
	}
	if dataStore == nil {
		coreLogger.Fatal("Failed to establish database connection")
		return nil
	}

	coreLogger.Infoln("**   proxy server utility Setup  **")
	var mapAPI *googlemapapi.GoogleMapAPI
	if configuration.UtilitySettings.GoogleMapAPIConfig != nil {
		if utility := googlemapapi.NewGoogleMapAPI(configuration.UtilitySettings.GoogleMapAPIConfig.Key); utility == nil {
			coreLogger.Fatal("Failed to set googlemapapi")
		} else {
			mapAPI = utility
		}
	}
	proxy := proxy.NewProxy(dataStore, mapAPI)
	if proxy == nil {
		coreLogger.Fatal("Proxy not created")
		return nil
	}

	servers := []types.Server{}
	if configuration.ServerSettings.APIServerConfig != nil {
		if apiServer := api.NewAPIServer(proxy, dataStore, mapAPI, configuration.ServerSettings.APIServerConfig.Port); apiServer == nil {
			coreLogger.Fatal("Failed to create apiserver")
		} else {
			servers = append(servers, apiServer)
		}
	}

	if configuration.ServerSettings.GraphqlServerConfig != nil {
		// TODO: graphql
		coreLogger.Infof("implementation of graphql is pending")
		// servers = append(servers, apiServer)
	}
	yayInstance = &Yay{
		config:  *configuration,
		servers: servers,
		proxy:   proxy,
	}
	return yayInstance
}

func (yay *Yay) GetConfig() types.Config {
	return yay.config
}

func (yay *Yay) Start() {
	coreLogger.Infoln("**          start server         **")
	for _, server := range yay.servers {
		if !server.Serve() {
			coreLogger.Error("Failed to serve")
		}
	}
}
