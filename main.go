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

	infoHandlerw()

	log.Fatal(http.ListenAndServe(":1010", router))
}
