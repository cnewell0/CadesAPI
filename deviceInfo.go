package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type info struct {
	DeviceID          string `json:"device_id"`
	UserAgent         string `json:"user_agent"`
	BatteryLevel      string `json:"battery_level"`
	ScreenOrientation string `json:"screen_orientation`
	IPAddress         string `json:"ip_address"`
}

//SomeInfo ... you know it's an array of the info's
type SomeInfo []info

var infoRecord = SomeInfo{
	{
		DeviceID:          "1234",
		UserAgent:         "RandomeUAios547",
		BatteryLevel:      "67",
		ScreenOrientation: "landscape",
		IPAddress:         "129.232.23.121",
	},
}

// InsertInfoRecord inserts a record from JSON into a fake db
func InsertInfoRecord(lilInfo info) {

	//fmt.Printf("%+v\n", anEvent)

	infoRecord = append(infoRecord, lilInfo)

	fmt.Printf("%+v\n", infoRecord)
}

// GetInfoRecord grabs and returns the all the events stored, from a mock db
func GetInfoRecord(deviceID string) SomeInfo {
	fmt.Println("Got 'em: Returning all events")
	//geoRecord = append(geoRecord, mdb.myEvents)

	return infoRecord
}

func infoHandlerw(router *mux.Router) {
	TheInfos := info{}

	router.Handle("/info", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Body read error, %v", err)
			w.WriteHeader(500) //500 Internal Server Error
			return
		}
		json.Unmarshal(reqBody, &TheInfos)
		InsertInfoRecord(TheInfos)
	})).Methods(http.MethodPost)

	router.Handle("/info/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		infoRecord = GetInfoRecord("{id}")
		json.NewEncoder(w).Encode(infoRecord)
	})).Methods(http.MethodGet)

}
