package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"github.com/jmoiron/sqlx"
	"quant/config"
)

// use sqlx to connect to the database
var global_db *sqlx.DB = nil

func GetGlobalDB() (*sqlx.DB, error) {
	if global_db != nil {
		return global_db, nil
	}
	// connect to the database
	// return the connection
	globalConfig := config.GetGlobalConfig()
	if globalConfig == nil {
		panic("globalConfig is nil")
	}
	db, err := GetDBConnection(globalConfig)
	if err != nil {
		return nil, err
	}
	global_db = db
	return global_db, nil
}

// get db connection by host, port, username, password, database
func GetDBConnection(config *config.Config) (*sqlx.DB, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}
	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Database)

	// Open a new database connection
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
