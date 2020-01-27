package storage

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/m1/go-service-healthcheck/config"
)

// MySQLDB ...
type SQLiteDB struct {
	*gorm.DB
}

// NewMySQLDB ...
func NewSQLiteDB(dbConfig config.DBConfig) (*SQLiteDB, error) {
	var err error
	db := &SQLiteDB{}
	db.DB, err = gorm.Open("sqlite3", dbConfig.File)
	db.LogMode(false)
	return db, err
}
