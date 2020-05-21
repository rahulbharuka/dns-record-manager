package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	// driver for mysql
	_ "github.com/go-sql-driver/mysql"
)

var (
	// initOnce protects the following
	initOnce  sync.Once
	dbHandler *handlerImpl
)

// Handler is a handler object for MySQL db with replicas.
type Handler interface {
	Query(query string) (*sql.Rows, error)
}

// handlerImpl ...
type handlerImpl struct {
	master *sql.DB
}

// NewDBHandler returns a db handler.
func NewDBHandler() Handler {
	initOnce.Do(func() {
		masterDB := openDB()
		dbHandler = &handlerImpl{
			master: masterDB,
		}
	})
	return dbHandler
}

// openDB opens master database.
func openDB() *sql.DB {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", username, password, host, dbName)
	masterDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panic(err.Error())
	}
	return masterDB
}

// Query loads data from table.
func (h *handlerImpl) Query(query string) (*sql.Rows, error) {
	return h.master.Query(query)
}
