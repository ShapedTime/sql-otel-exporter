package dbreader

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/sijms/go-ora/v2"
	"log"
	otelhelpers "sqlquery_otel_exporter/otel-helpers"
)

type DbReader struct {
	db                   *sql.DB
	sqlQuery             string
	trackingColumn       string
	defaultTrackingValue string
	columnTracker        *ColumnTracker
	name                 string
}

func NewDbReader(sqlQuery, trackingColumn, defaultTrackingValue, driver, dataSource string, columnTracker *ColumnTracker, name string) (*DbReader, error) {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		return nil, err
	}

	return &DbReader{
		db:                   db,
		sqlQuery:             sqlQuery,
		trackingColumn:       trackingColumn,
		defaultTrackingValue: defaultTrackingValue,
		columnTracker:        columnTracker,
		name:                 name,
	}, nil
}

func (r *DbReader) Read() ([]map[string]string, error) {
	args := []interface{}{}

	// If we have a tracking column, we need to add the tracking value to the query
	if r.trackingColumn != "" {
		trackingValue := r.columnTracker.GetTrackingColumn(r.trackingColumn)

		if trackingValue == "" {
			trackingValue = r.defaultTrackingValue
		}
		args = []interface{}{trackingValue}
	}

	query := r.sqlQuery

	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]string
	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		row := make(map[string]string)
		for i, colName := range cols {
			val, ok := columnPointers[i].(*interface{})
			if ok {
				row[colName] = fmt.Sprintf("%v", *val)
			}
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(results) > 0 {
		lastRow := results[len(results)-1]

		// If we have a tracking column, we need to update the tracking value
		if r.trackingColumn != "" {
			lastTrackingValue := lastRow[r.trackingColumn]
			if err := r.columnTracker.SetTrackingColumn(r.trackingColumn, lastTrackingValue); err != nil {
				return nil, err
			}
		}
	}

	otelhelpers.IncReceivedRowsCounter(int64(len(results)), r.name)

	return results, nil
}

func (r *DbReader) Close() {
	log.Println("Closing db connection")
	r.db.Close()
}
