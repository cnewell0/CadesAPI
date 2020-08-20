package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func initLog() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	ExtraGeoRecord := []event{}
	GeoRecord := event{}

	router := mux.NewRouter().StrictSlash(true)
	dba, _ := Configure("fake", "okay", "random", "fillers")

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

	infoHandlerw(router)
	HandleTime(router)

	var w http.ResponseWriter
	pullMsgsConcurrenyControl(w, "kochava-testing", "test-sub")

	// InfoLogger.Println("Starting the application...")
	// InfoLogger.Println("Something noteworthy happened")
	// WarningLogger.Println("There is something you should know about")
	// ErrorLogger.Println("Something went wrong")

	log.Fatal(http.ListenAndServe(":1010", router))
}
