package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type event struct {
	DeviceId  string `json:"device_id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	IpAddress string `json:"ip_address"`
}

type info struct {
	DeviceId          string `json:"device_id"`
	UserAgent         string `json:"user_agent"`
	BatteryLevel      string `json:"battery_level"`
	IpAddress         string `json:"ip_address"`
	ScreenOrientation string `json: "screen_orientation`
}

type events []event

var allEvents = events{
	{
		DeviceId:  "1234-123919291-123-12312",
		Latitude:  "48.121",
		Longitude: "127.12",
		IpAddress: "129.232.23.121",
	},
}

type allInfo []info

var infos = allInfo{}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	//myRouter.HandleFunc("/", homeLink)

	myRouter.HandleFunc("/geo", createNewEvent).Methods("POST")
	myRouter.HandleFunc("/geo", retrunAllEvents).Methods("GET")
	myRouter.HandleFunc("/geo/{deviceId}", returnSingleEvent).Methods("GET")
	myRouter.HandleFunc("/geo/{deviceId}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/geo/{deviceId}", updateEvent).Methods("PATCH")

	log.Fatal(http.ListenAndServe(":1010", myRouter))
}

func createNewEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event

	if r.Method != http.MethodPost {
		w.WriteHeader(405) //405 Mehod Not Allowed
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Body read error, %v", err)
		w.WriteHeader(500) //500 Internal Server Error
		return
	}

	if err = json.Unmarshal(reqBody, &newEvent); err != nil {
		log.Printf("Body parse error, %v", err)
		w.WriteHeader(400) //400 Bad Request
		return
	}

	if newEvent.DeviceId == "" || newEvent.Latitude == "" || newEvent.Longitude == "" || newEvent.IpAddress == "" {
		log.Printf("incorrect keys passed")
		w.WriteHeader(400) //400 Bad Request
		return
	}

	allEvents = append(allEvents, newEvent)
	w.WriteHeader(http.StatusCreated)

	fmt.Printf("%+v\n", allEvents)
	json.NewEncoder(w).Encode(newEvent)

}

// func homeLink(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome home!")
// 	fmt.Println("Endpoint Hit: beep boop")
// }

func retrunAllEvents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got 'em: Returning all events")
	json.NewEncoder(w).Encode(allEvents)
}

func returnSingleEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["deviceId"]

	for _, singleEvent := range allEvents {
		if singleEvent.DeviceId == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["deviceId"]

	for i, singleEvent := range allEvents {
		if singleEvent.DeviceId == eventID {
			allEvents = append(allEvents[:i], allEvents[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Enter in proper JSON formatting")
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range allEvents {
		if singleEvent.DeviceId == mux.Vars(r)["deviceId"] {
			singleEvent.Latitude = updatedEvent.Latitude
			singleEvent.Longitude = updatedEvent.Longitude
			allEvents = append(allEvents[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func main() {
	handleRequests()
}
