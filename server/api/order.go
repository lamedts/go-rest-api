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
		outgoingJSON, apiError = server.updateOrder(mux.Vars(r)["id"], r.Body)
	}
	if apiError.Code != 0 {
		w.WriteHeader(apiError.Code)
		fmt.Fprint(w, fmt.Sprintf(`{"error": "%v"}`, apiError.Error))
	}
	fmt.Fprint(w, string(outgoingJSON))
}

// listOrder is to list all order according some param
func (server *APIServer) listOrder(page int, limit int) ([]byte, yayerror.APIError) {
	customAPI500 := yayerror.API500
	if orders, err := server.db.ReadOrder(page, limit); err != nil {
		APILogger.Warnf("Error when listing: %+v", err)
		customAPI500.Error = err.Error()
		return nil, customAPI500
	} else if orders != nil {
		type orderJSON = struct {
			ID       int     `json:"id"`
			Distance float32 `json:"distance"`
			Status   string  `json:"status"`
		}
		var ordersJSON []orderJSON
		for _, order := range *orders {
			tmpJSON := orderJSON{
				ID:       order.ID,
				Distance: order.Distance,
				Status:   order.Status,
			}
			ordersJSON = append(ordersJSON, tmpJSON)
		}
		returnJSON, _ := json.Marshal(ordersJSON)
		return []byte(returnJSON), yayerror.APIError{}
	}
	APILogger.Warnf("Unknow Error when doing listOrder request")
	customAPI500.Error = "Unknow Error"
	return nil, customAPI500
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

	if len(locationData.Destination) > 0 && len(locationData.Origin) > 0 {

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
			APILogger.Warnf("Cant parse a request data: %+v", locationData)
			customAPI500.Error = "Parse Float error"
			return nil, customAPI500
		}
		// google api
		order.Distance = 19

		if order, err := server.db.CreateOrder(order); err != nil {
			APILogger.Warnf("Eror when creating: %+v", err)
			customAPI500.Error = err.Error()
			return nil, customAPI500
		} else if order != nil {
			APILogger.Infof("Order created, id: %+v", order.ID)
			type returnJSON = struct {
				ID       int     `json:"id"`
				Distance float32 `json:"distance"`
				Status   string  `json:"status"`
			}
			createdOrder, _ := json.Marshal(returnJSON{
				ID:       order.ID,
				Status:   order.Status,
				Distance: order.Distance,
			})
			return []byte(createdOrder), yayerror.APIError{}
		}
	} else {
		APILogger.Warnf("Bad Request Data: %+v", locationData)
		customAPI500.Error = "Bad Request Data"
		return nil, customAPI500
	}

	APILogger.Warnf("Unknow Error when doing creatOrder request")
	customAPI500.Error = "Unknow Error"
	return nil, customAPI500
}

// updateOrder is to update the info and take order
// only one field will be updated at this moment,
// ** change the state of the order
func (server *APIServer) updateOrder(orderIDstr string, requestData io.ReadCloser) ([]byte, yayerror.APIError) {
	customAPI500 := yayerror.API500
	type incomingData struct {
		Status string `json:"status"`
	}
	var statusData incomingData
	var orderID int
	json.NewDecoder(requestData).Decode(&statusData)

	if orderIDstr == "" {
		APILogger.Warnf("Bad Request order id: %+v", statusData)
		customAPI500.Error = "Bad Request order id"
		return nil, customAPI500
	} else {
		if result, err := strconv.Atoi(orderIDstr); err != nil {
			APILogger.Warnf("Cant parse a orderID param: %+v", orderIDstr)
			customAPI500.Error = "Bad orderID"
			return nil, customAPI500
		} else {
			orderID = result
		}
	}

	if statusData.Status == "taken" {
		// key of the updateData must be a col in the table
		// value must be in correct data type
		updateData := map[string]interface{}{"status": "ASSIGNED"}
		if isTaken, err := server.db.UpdateOrder(orderID, updateData); err != nil {
			APILogger.Warnf("Eror when updating: %+v", err)
			customAPI500.Error = err.Error()
			return nil, customAPI500
		} else if !isTaken {
			APILogger.Infof("Order taken, id: %+v", orderID)
			return []byte(`{"status": "SUCCESS"}`), yayerror.APIError{}
		}
		return nil, yayerror.API409

	} else {
		APILogger.Warnf("Bad Request Data: %+v", statusData)
		customAPI500.Error = "Bad Request Data"
		return nil, customAPI500
	}

	APILogger.Warnf("Unknow Error when doing updateOrder request")
	customAPI500.Error = "Unknow Error"
	return nil, customAPI500
}
