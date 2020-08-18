package main

import "fmt"

// MockSQLAccessor is a struct for a fake db to test on
type MockSQLAccessor struct {
	hostname string
	myEvents []event
}

// NewMockSQLAccessor is the contstuctor for this fake db
func NewMockSQLAccessor(hostname string, myEvents []event) *MockSQLAccessor {
	mockSQL := MockSQLAccessor{
		hostname: hostname,
		myEvents: myEvents,
	}
	return &mockSQL
}

// InsertGeoRecord inserts a record from JSON into a fake db
func (mdb MockSQLAccessor) InsertGeoRecord(anEvent event) {

	//fmt.Printf("%+v\n", anEvent)

	geoRecord = append(geoRecord, anEvent)

	fmt.Printf("%+v\n", geoRecord)
}

// GetGeoRecord grabs and returns the all the events stored, from a mock db
func (mdb MockSQLAccessor) GetGeoRecord(deviceID string) SomeEvent {
	fmt.Println("Got 'em: Returning all events")
	//geoRecord = append(geoRecord, mdb.myEvents)

	return geoRecord
}
