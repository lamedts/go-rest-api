package datastore

import (
	"errors"
	"fmt"
	"go-rest-api/types"
	"go-rest-api/types/datastore"
	"strings"
	"sync"

	"go-rest-api/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var databaseLogger = logger.GetLogger("datastore")

type mysqlDB struct {
	*sqlx.DB
	mutex *sync.RWMutex
}

func NewMysqlDB(mysqlConfig types.DatastoreSettings) *mysqlDB {
	mysqlConnectionString := fmt.Sprintf("%s:%s@(%s:%d)/%s", mysqlConfig.User, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.DBname)
	sqlxdb, err := sqlx.Connect("mysql", mysqlConnectionString)
	if err != nil {
		databaseLogger.Fatalln("Failed to connect to database:", err)
	}

	if err := sqlxdb.Ping(); err != nil {
		databaseLogger.Fatal(err)
		return nil
	}

	databaseLogger.Infoln("database connection is established")

	mysqlDB := mysqlDB{DB: sqlxdb, mutex: &sync.RWMutex{}}

	// TODO: smiple DataMigration

	return &mysqlDB
}

/**
 * order CRUD
 */
func (db *mysqlDB) ReadOrder(page int, limit int) (bool, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// if _, err := db.NamedExec(`
	// 	INSERT INTO settlement (settlement_date, details) VALUES (:settlement_date, :details)`, settlementData); err != nil {
	// 	databaseLogger.WithFields(logrus.Fields{
	// 		"Flow": "datastore",
	// 		"func": "SaveSettlement",
	// 	}).Warn(err)
	// 	return false
	// }
	return true, nil
}

// UpdateOrder
// return (bool for isTaken, error)
func (db *mysqlDB) UpdateOrder(orderID int, updateData map[string]interface{}) (bool, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	var (
		cols   = ""
		values = ""
	)
	for key := range updateData {
		cols = cols + key + ","
		values = values + updateData[key].(string) + ","
	}
	updateStmt := fmt.Sprintf(`UPDATE orders set %s = "%s" WHERE id = :order_id`, strings.TrimSuffix(cols, ","), strings.TrimSuffix(values, ","))
	if result, err := db.NamedExec(updateStmt, map[string]interface{}{"order_id": orderID}); err != nil {
		databaseLogger.WithFields(logrus.Fields{
			"Flow": "datastore",
			"func": "UpdateOrder",
		}).Warn(err)
		return false, err
	} else if rowNum, _ := result.RowsAffected(); rowNum == 0 {
		return true, nil
	} else if rowNum == 1 {
		return false, nil
	}
	return false, errors.New("Unknow")
}

func (db *mysqlDB) CreateOrder(order datastore.Order) (*datastore.Order, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	tx, _ := db.Beginx()

	var insertedDestID int64
	var insertOriginID int64
	var insertOrderID int64
	if result, err := tx.NamedExec(`INSERT INTO orders_origin (latitude,longtitude) VALUES (:latitude,:longtitude);`, &order.OrderOrigin); err != nil {
		databaseLogger.WithFields(logrus.Fields{
			"Flow":   "datastore",
			"func":   "CreateOrder",
			"Action": "get OrderOrigin",
		}).Error(err)
		tx.Rollback()
		return nil, err
	} else {
		insertedDestID, _ = result.LastInsertId()
	}

	if result, err := tx.NamedExec(`INSERT INTO orders_destination (latitude,longtitude) VALUES (:latitude,:longtitude);`, &order.OrderDestination); err != nil {
		databaseLogger.WithFields(logrus.Fields{
			"Flow":   "datastore",
			"func":   "CreateOrder",
			"Action": "get OrderDestination",
		}).Error(err)
		tx.Rollback()
		return nil, err
	} else {
		insertOriginID, _ = result.LastInsertId()
	}

	// insert order info
	queryInterface := map[string]interface{}{
		"origin_id":      insertOriginID,
		"destination_id": insertedDestID,
		"distance":       order.Distance,
	}
	if result, err := tx.NamedExec(`INSERT INTO orders (origin_id,destination_id,distance) VALUES (:origin_id,:destination_id,:distance)`, queryInterface); err != nil {
		databaseLogger.WithFields(logrus.Fields{
			"Flow":   "datastore",
			"func":   "CreateOrder",
			"Action": "inser order",
		}).Error(err)
		tx.Rollback()
		return nil, err
	} else {
		insertOrderID, _ = result.LastInsertId()
	}

	if orderStmt, err := tx.PrepareNamed(`SELECT id, status FROM orders WHERE id = :order_id`); err != nil {
		databaseLogger.WithFields(logrus.Fields{
			"Flow":   "datastore",
			"func":   "CreateOrder",
			"Action": "prepare get order",
		}).Error(err)
		tx.Rollback()
		return nil, err
	} else {
		var fetchedOrder datastore.Order
		if err := orderStmt.Get(&fetchedOrder, map[string]interface{}{"order_id": insertOrderID}); err != nil {
			databaseLogger.WithFields(logrus.Fields{
				"Flow":   "datastore",
				"func":   "CreateOrder",
				"Action": "Get Order",
			}).Error(err)
			tx.Rollback()
			return nil, err
		} else if err := tx.Commit(); err != nil {
			databaseLogger.WithFields(logrus.Fields{
				"Flow":   "datastore",
				"func":   "CreateOrder",
				"Action": "commit",
			}).Error(err)
			tx.Rollback()
			return nil, err
		}
		order.Status = fetchedOrder.Status
		order.ID = fetchedOrder.ID
		order.CreatedTime = fetchedOrder.CreatedTime
		order.UpdatedTime = fetchedOrder.UpdatedTime
		databaseLogger.Debugf("%+v", order)
		return &order, nil
	}
}
