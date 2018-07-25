package api_test

import (
	"go-rest-api/server/api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPathOrderHandler(t *testing.T) {
	// check route
	req, _ := http.NewRequest("GET", "/hello/chris", nil)
	res := httptest.NewRecorder()

	server := api.NewAPIServer(nil, nil, nil, 8080)
	server.OrderHandler(res, req)

	if res.Body.String() != `{"code":400,"error":"","key":"NOT_FOUND"}` {
		t.Error(`Fail! It should {"code":400,"error":"","key":"NOT_FOUND"}"`)
	}
}

func TestPutOrderHandler(t *testing.T) {
	// check route
	req, _ := http.NewRequest("PUT", "/orders/123", strings.NewReader(""))
	res := httptest.NewRecorder()

	server := api.NewAPIServer(nil, nil, nil, 8080)
	server.OrderHandler(res, req)

	t.Log(res.Body)
	if res.Body.String() != `{"error": "Bad Request order id"}` {
		t.Error(`Fail! It should {"error": "Bad Request order id"}"`)
	}
}
