package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func databaseConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.Username, config.Password, config.Hostname, config.Port, config.Database)
}

// Setups the connection to the datbase, then checks if it was successful via the ping
func connection() (*sql.DB, error) {
	dbTemp, err := sql.Open("mysql", databaseConnectionString())
	if err != nil {
		log.Error("Could not connect to Database: ", err.Error())
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = dbTemp.PingContext(ctx)
	if err != nil {
		log.Error("Could not connect to database: ping failed: ", err.Error())
		return nil, err
	}
	log.Debug(fmt.Sprintf("Successfully connected to database `%s`", config.Database))
	return dbTemp, nil
}

// Executes a ping to the database in order to check if it is online
func CheckDatabase() error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	if err := db.PingContext(ctx); err != nil {
		log.Error("Database health check failed: ", err.Error())
		return err
	}
	return nil
}
