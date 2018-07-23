package errors

type APIError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Key   string `json:"key"`
}

var (
	API400 = APIError{Code: 400, Key: "NOT_FOUND"}
	API409 = APIError{Code: 409, Error: "ORDER_ALREADY_BEEN_TAKEN"}
	API500 = APIError{Code: 500}
)
