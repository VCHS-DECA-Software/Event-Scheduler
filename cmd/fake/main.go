package main

import (
	"Event-Scheduler/output"
	"Event-Scheduler/scheduler"
	"log"
	"os"

	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	gofakeit.Seed(0)

	c := FakeContext(ContextOptions{
		Students:      100,
		Judges:        25,
		Rooms:         5,
		TimeDivisions: 6,
		//* uncomment this for a more comprehendible result
		// Students:      10,
		// Judges:        3,
		// Rooms:         1,
		// TimeDivisions: 5,

		Events:        10,
		GroupCapacity: 3,
	})
	requests := FakeRequests(c, RequestOptions{
		MinimumEvents: 1,
		MaximumEvents: 3,
		SoloRatio:     0.2,
	})

	o := scheduler.Schedule(c, requests)

	f, err := os.Create("output.csv")
	if err == os.ErrExist {
		f, err = os.Open("output.csv")
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	err = output.CSV(f, o)
	if err != nil {
		log.Fatal(err)
	}
}
