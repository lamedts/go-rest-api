package datastore

import (
	"time"
)

// OrderStatusType: two types
const (
	OrderStatusAssigned string = "ASSIGNED"
	OrderStatusUnassign string = "UNASSIGN "
)

// Order contains all info
type Order struct {
	ID               int       `db:"id" json:"order_id"`
	Status           string    `db:"status" json:"status"`
	UpdatedTime      time.Time `db:"updated_at" json:"updated_at"`
	CreatedTime      time.Time `db:"created_at" json:"created_at"`
	OrderDestination *OrderDestination
	OrderOrigin      *OrderOrigin
}

// OrderDestination contains info of destination only
type OrderDestination struct {
	ID          int       `db:"id" json:"order_destination_id"`
	Latitude    float32   `db:"latitude" json:"latitude"`
	Longtitude  float32   `db:"longtitude" json:"longtitude"`
	UpdatedTime time.Time `db:"updated_at" json:"updated_at"`
	CreatedTime time.Time `db:"created_at" json:"created_at"`
}

// OrderOrigin contains info of origin only
type OrderOrigin struct {
	ID          int       `db:"id" json:"order_origin_id"`
	Latitude    float32   `db:"latitude" json:"latitude"`
	Longtitude  float32   `db:"longtitude" json:"longtitude"`
	UpdatedTime time.Time `db:"updated_at" json:"updated_at"`
	CreatedTime time.Time `db:"created_at" json:"created_at"`
}
