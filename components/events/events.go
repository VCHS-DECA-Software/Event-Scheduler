package events

import (
	"main/components/dbmanager"
)

type Event struct {
	ID          string `storm:"id"`
	AdminID     string `storm:"index"`
	Name        string `storm:"unique"`
	Description string
	EventType   string `storm:"index"`
	StartTime   string
	EndTime     string
	Location    string `storm:"index"`
	MaxStudents int
	JudgeIDs    []string
	TeamIDs     []string
}

func GetEvent(id string) (Event, error) {
	var event Event
	err := dbmanager.Query("ID", id, &event)
	return event, err
}

//GetAllEvents is for internal use only. Use ViewAllEvents from the users package to view all events.
func GetAllEvents() ([]Event, error) {
	var events []Event
	err := dbmanager.QueryAll(&events)
	return events, err
}
