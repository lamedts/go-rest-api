package api

import (
	"encoding/json"
	"fmt"
	yayerror "go-rest-api/errors"
	"go-rest-api/types/datastore"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (server *APIServer) OrderHandler(w http.ResponseWriter, r *http.Request) {

	var outgoingJSON, err = json.Marshal(yayerror.API400)
	var apiError = yayerror.APIError{}
	if err != nil {
		APILogger.Errorf("Cant marshal json: %+v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if mux.Vars(r)["page"] != "" && mux.Vars(r)["limit"] != "" {
		page, atoiErr1 := strconv.Atoi(mux.Vars(r)["page"])
		limit, atoiErr2 := strconv.Atoi(mux.Vars(r)["limit"])
		if atoiErr1 == nil && atoiErr2 == nil {
			outgoingJSON, apiError = server.listOrder(page, limit)
		}
	} else if r.Method == "POST" {
		outgoingJSON, apiError = server.createOrder(r.Body)
	} else if r.Method == "PUT" {
		outgoingJSON, apiError = server.updateOrder()
	}
	if apiError.Code != 0 {
		w.WriteHeader(apiError.Code)
		fmt.Fprint(w, fmt.Sprintf(`{"error": "%v"}`, apiError.Error))
	}
	fmt.Fprint(w, string(outgoingJSON))
}

// listOrder is to list all order according some param
func (server *APIServer) listOrder(page int, limit int) ([]byte, yayerror.APIError) {
	return []byte(`{"action":"listOrder"}`), yayerror.APIError{}
}

// createOrder is to place order
func (server *APIServer) createOrder(requestData io.ReadCloser) ([]byte, yayerror.APIError) {
	customAPI500 := yayerror.API500
	type incomingData struct {
		Origin      []string `json:"origin"`
		Destination []string `json:"destination"`
	}
	var locationData incomingData
	json.NewDecoder(requestData).Decode(&locationData)

	destLat, parseErr1 := strconv.ParseFloat(locationData.Destination[0], 32)
	destLong, parseErr2 := strconv.ParseFloat(locationData.Destination[1], 32)
	origLat, parseErr3 := strconv.ParseFloat(locationData.Origin[0], 32)
	origLong, parseErr4 := strconv.ParseFloat(locationData.Origin[1], 32)

	order := datastore.Order{
		OrderDestination: &datastore.OrderDestination{
			Latitude:   float32(destLat),
			Longtitude: float32(destLong),
		},
		OrderOrigin: &datastore.OrderOrigin{
			Latitude:   float32(origLat),
			Longtitude: float32(origLong),
		},
	}
	if parseErr1 != nil || parseErr2 != nil || parseErr3 != nil || parseErr4 != nil {
		APILogger.Warnf("Cant parse a data from a creatOrder request")
		customAPI500.Error = "Parse Float error"
		return nil, customAPI500
	}

	if order, _ := server.db.CreateOrder(order); order == nil {
		return []byte("createOrder!\n"), yayerror.APIError{}
	} else {

	}
	APILogger.Warnf("Unknow Erro when doing creatOrder request")
	customAPI500.Error = "Unknow Error"
	return nil, customAPI500
}

// updateOrder is to update the info and take order
func (server *APIServer) updateOrder() ([]byte, yayerror.APIError) {
	return []byte("updateOrder!\n"), yayerror.APIError{}
}
