package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

//DBAccess Interface to keep track of CRUD methods
type DBAccess interface {
	InsertGeoRecord(anEvent event)
	GetGeoRecord(deviceID string) SomeEvent
}

// MockSQLAccessor is a struct for a fake db to test on
type MockSQLAccessor struct {
	hostname string
	myEvents []event
}

// SQLAcessor is a struct for connecting to the actual db
type SQLAcessor struct {
	dbRead  *sql.DB
	anEvent []event
}

type event struct {
	DeviceID  string `json:"device_id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	IPAddress string `json:"ip_address"`
}

// SomeEvent is, you know, some event
type SomeEvent []event

var geoRecord = SomeEvent{
	{
		DeviceID:  "1234",
		Latitude:  "50.111",
		Longitude: "100.222",
		IPAddress: "129.232.23.121",
	},
}

// NewMockSQLAccessor is the contstuctor for this fake db
func NewMockSQLAccessor(hostname string, myEvents []event) *MockSQLAccessor {
	mockSQL := MockSQLAccessor{
		hostname: hostname,
		myEvents: myEvents,
	}
	return &mockSQL
}

// NewSQLAccessor is the constructor for mysql
func NewSQLAccessor(dbRead *sql.DB, anEvent []event) *SQLAcessor {
	dbGlobal := SQLAcessor{
		dbRead:  dbRead,
		anEvent: anEvent,
	}
	return &dbGlobal
}

func configure(hostname string, port string, username string, password string) (DBAccess, error) {
	var dba DBAccess
	if hostname == "fake" {
		dba = NewMockSQLAccessor(hostname, geoRecord)
	} else {
		configRead := mysql.Config{
			User:   username,
			Passwd: password,
			Net:    "tcp",
			Addr:   fmt.Sprintf("%s:%d", hostname, port),
		}
		dbRead, err := sql.Open("mysql", configRead.FormatDSN())
		if err != nil {
			return nil, errors.Wrap(err, "failed to open db global read client")
		}
		if err := dbRead.Ping(); err != nil {
			return nil, errors.Wrap(err, "failed to ping db global read server")
		}
		dba = NewSQLAccessor(dbRead, geoRecord)
	}
	return dba, nil
}

// InsertGeoRecord inserts a record from JSON into a fake db
func (mdb MockSQLAccessor) InsertGeoRecord(anEvent event) {

	//fmt.Printf("%+v\n", anEvent)

	geoRecord = append(geoRecord, anEvent)

	fmt.Printf("%+v\n", geoRecord)
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

// GetGeoRecord grabs and returns the all the events stored, from a mock db
func (mdb MockSQLAccessor) GetGeoRecord(deviceID string) SomeEvent {
	fmt.Println("Got 'em: Returning all events")
	//geoRecord = append(geoRecord, mdb.myEvents)

	return geoRecord
}

func main() {
	ExtraGeoRecord := []event{}
	GeoRecord := event{}

	router := mux.NewRouter().StrictSlash(true)
	dba, _ := configure("fake", "okay", "random", "fillers")

	router.Handle("/geo/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ExtraGeoRecord = dba.GetGeoRecord("{id}")
		ExtraGeoRecord = append(ExtraGeoRecord, GeoRecord)
		json.NewEncoder(w).Encode(ExtraGeoRecord)
	})).Methods(http.MethodGet)

	router.Handle("/geo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Body read error, %v", err)
			w.WriteHeader(500) //500 Internal Server Error
			return
		}
		json.Unmarshal(reqBody, &GeoRecord)
		dba.InsertGeoRecord(GeoRecord)
	})).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":1010", router))
}
