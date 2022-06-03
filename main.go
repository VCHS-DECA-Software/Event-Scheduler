package main

import (
	"main/components/dbmanager"
	"main/components/events"
	"main/components/users"
)

func main() {
	err := dbmanager.Open("eventscheduler.db")
	if err != nil {
		panic(err)
	}

	var errors []error

	errors = append(errors, dbmanager.AutoCreateStruct(&events.Event{}))
	errors = append(errors, dbmanager.AutoCreateStruct(&users.Admin{}))
	errors = append(errors, dbmanager.AutoCreateStruct(&users.Judge{}))
	errors = append(errors, dbmanager.AutoCreateStruct(&users.Student{}))
	errors = append(errors, dbmanager.AutoCreateStruct(&users.Team{}))
	checkForErrors(errors)
	
	err := dbmanager.Close("eventscheduler.db")
	if err != nil {
		panic(err)
	}

}

func checkForErrors(errors []error) {
	for _, err := range errors {
		if err != nil {
			panic(err)
		}
	}
}
