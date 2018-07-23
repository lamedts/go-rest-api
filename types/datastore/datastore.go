package datastore

// DatastoreType
const (
	DatastoreTypeMysql   string = "mysql"
	DatastoreTypeMongodb string = "mongodb"
)

type OrderDataStore interface {
	UpdateOrder(orderID int, updateData map[string]interface{}) (bool, error)
	ReadOrder(page int, limit int) (*[]Order, error)
	CreateOrder(order Order) (*Order, error)
}

type DataStore interface {
	OrderDataStore
}
