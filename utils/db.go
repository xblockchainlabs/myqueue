package utils

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"
)

type Connection struct {
	User    string
	Pass    string
	Host    string
	Port    string
	Name    string
	LogMode bool
}

var dbConnection *Connection
var dbObj *gorm.DB
var dbOnce sync.Once

func (c *Connection) isSet() bool {
	return len(c.Host) > 0 && len(c.Name) > 0
}

func (c *Connection) getDSN() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.Name)
}

func SetConnection(user string, pass string, host string, port string, name string, log bool) {
	dbConnection = &Connection{
		user,
		pass,
		host,
		port,
		name,
		log,
	}
}

// Opening a database and save the reference to `Database` struct.
func newDB() (db *gorm.DB, err error) {
	dsn := dbConnection.getDSN()
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}
	db.DB().SetMaxIdleConns(10)
	db.LogMode(dbConnection.LogMode)
	return
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() (db *gorm.DB, err error) {
	if dbConnection == nil || !dbConnection.isSet() {
		err = errors.New("DB connection is not set")
		return
	}

	dbOnce.Do(func() {
		db, err = newDB()
		if err == nil {
			dbObj = db
		}
	})

	db = dbObj
	return
}
