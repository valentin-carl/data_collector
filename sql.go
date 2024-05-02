package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func create(dbc *sql.Conn, tableName string) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (timestamp INTEGER, total_workflow_duration REAL);", tableName)
	log.Printf("create query: '%s'", query)
	_, err := dbc.ExecContext(context.Background(), query)
	log.Printf("created table %s or already exists", tableName)
	return err
}

func insertMeasurement(dbc *sql.Conn, tableName string, timeStamp uint64, totalWorkflowDuration float64) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (timestamp, total_workflow_duration) VALUES (%d, %f);", tableName, timeStamp, totalWorkflowDuration)
	log.Printf("insert query: '%s'", query)
	res, err := dbc.ExecContext(context.Background(), query)
	if err != nil {
		log.Printf("error occurred while trying to insert values into DB: '%s'", err.Error())
		return -1, err
	}
	rA, err := res.RowsAffected()
	if err != nil {
		log.Printf("could not get number of rows affected by insert query: '%s'", err.Error())
		return -1, err
	}
	log.Printf("inserting values successful, rows affected: %d", rA)
	return rA, nil
}
