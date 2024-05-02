package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
)

/*
database scheme:
- timestamp INTEGER
- total_workflow_duration REAL

assuming request structure
{
  "tableName": "test_table_1234    // to differentiate measurement runs
  "timestamp": 1234, 			   // unix milli
  "totalWorkflowDuration": 12.34   // in milliseconds
}
*/

const connectionString = "file:database.db"
const listenAddress = ":80"

func main() {

	ctx := context.Background()

	db, err := sql.Open("sqlite3", connectionString)
	handle(err)
	defer db.Close()

	dbc, err := db.Conn(ctx)
	defer dbc.Close()
	handle(err)

	err = dbc.PingContext(ctx)
	handle(err)

	log.Fatal(serve(listenAddress, dbc))
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

type Measurement struct {
	TableName             string  `json:"tableName"`
	Timestamp             uint64  `json:"timestamp"`
	TotalWorkflowDuration float64 `json:"totalWorkflowDuration"`
}

func serve(address string, dbc *sql.Conn) error {

	e := func(err error, status int, w http.ResponseWriter) {
		if err != nil {
			log.Printf("got error: '%s'", err.Error())
			w.WriteHeader(status)
			w.Write([]byte(fmt.Sprintf("got error: '%s'", err.Error())))
		}
	}

	http.HandleFunc("/insert", func(w http.ResponseWriter, r *http.Request) {

		b, err := io.ReadAll(r.Body)
		e(err, http.StatusBadRequest, w)
		log.Printf("got request body: %s", string(b))

		var m Measurement
		err = json.NewDecoder(bytes.NewReader(b)).Decode(&m)
		e(err, http.StatusBadRequest, w)

		// create table if it doesn't exist yet
		err = create(dbc, m.TableName)
		e(err, http.StatusInternalServerError, w)

		rA, err := insertMeasurement(dbc, m.TableName, m.Timestamp, m.TotalWorkflowDuration)
		e(err, http.StatusInternalServerError, w)
		log.Printf("rows affected: %d", rA)

		w.WriteHeader(http.StatusAccepted)
	})
	return http.ListenAndServe(address, nil)
}
