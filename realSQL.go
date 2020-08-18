package main

import (
	"database/sql"
	"fmt"
)

// SQLAcessor is a struct for connecting to the actual db
type SQLAcessor struct {
	dbRead  *sql.DB
	anEvent []event
}

// NewSQLAccessor is the constructor for mysql
func NewSQLAccessor(dbRead *sql.DB, anEvent []event) *SQLAcessor {
	dbGlobal := SQLAcessor{
		dbRead:  dbRead,
		anEvent: anEvent,
	}
	return &dbGlobal
}

// InsertGeoRecord inserts a record from JSON into real mysql
func (rsql SQLAcessor) InsertGeoRecord(anEvent event) {
	stmt, err := rsql.dbRead.Prepare("INSERT INTO posts(deviceID) VALUES(?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	deviceID := "device_id"

	_, err = stmt.Exec(deviceID)
	if err != nil {
		panic(err.Error())
	}
	fmt.Print("New post was created")
}

// GetGeoRecord gets the record and returns it to the console
func (rsql SQLAcessor) GetGeoRecord(deviceID string) SomeEvent {

	result, err := rsql.dbRead.Query("SELECT device_id from posts WHERE device_id = ?")
	if err != nil {
		panic(err.Error())
	}

	var myEvents event

	defer result.Close()
	for result.Next() {
		err := result.Scan(&myEvents.DeviceID, &myEvents.Latitude, &myEvents.Longitude, &myEvents.IPAddress)
		if err != nil {
			panic(err.Error())
		}
	}

	geoRecord = append(geoRecord, myEvents)
	return geoRecord
}
