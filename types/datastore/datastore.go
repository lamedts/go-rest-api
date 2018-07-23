package datastore

// DatastoreType
const (
	DatastoreTypeMysql   string = "mysql"
	DatastoreTypeMongodb string = "mongodb"
)

type OrderDataStore interface {
	UpdateOrder(orderID int) (bool, error)
	ReadOrder(page int, limit int) (bool, error)
	CreateOrder(order Order) (*Order, error)
}

type DataStore interface {
	OrderDataStore
}
