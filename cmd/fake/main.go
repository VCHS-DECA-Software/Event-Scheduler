package main

import (
	"Event-Scheduler/output"
	"Event-Scheduler/scheduler"
	"flag"
	"os"

	"github.com/brianvoe/gofakeit/v6"
)

func main() {
	gofakeit.Seed(0)

	students := flag.Int("students", 100, "the number of students to generate")
	judges := flag.Int("judges", 25, "the number of judges to generate")
	rooms := flag.Int("rooms", 5, "the number of rooms to generate")
	divisions := flag.Int("divisions", 6, "the number of time slots to use")

	events := flag.Int("events", 4, "the number of different events to generate")
	roomCapacity := flag.Int("room-capacity", 4, "the maximum judges a room can hold")
	groupCapacity := flag.Int("group-capacity", 3, "the maximum capacity of a group")
	judgeTalent := flag.Int("judge-talent", 2, "the maximum number of different events a given judge can judge")

	minEvents := flag.Int("min-events", 1, "the minimum number of events a student will join")
	maxEvents := flag.Int("max-events", 2, "the maximum number of events a student will join")

	flag.Parse()

	c := FakeContext(ContextOptions{
		Students:      *students,
		Judges:        *judges,
		Rooms:         *rooms,
		TimeDivisions: *divisions,

		Events:         *events,
		GroupCapacity:  *groupCapacity,
		RoomCapacity:   *roomCapacity,
		MaxJudgeTalent: *judgeTalent,
	})
	requests := FakeRequests(c, RequestOptions{
		MinimumEvents: *minEvents,
		MaximumEvents: *maxEvents,
		SoloRatio:     0.2,
	})

	o := scheduler.Schedule(c, requests)

	f, err := os.Create("output.csv")
	if err != nil {
		panic(err)
	}

	err = output.CSV(f, o)
	if err != nil {
		panic(err)
	}
}
