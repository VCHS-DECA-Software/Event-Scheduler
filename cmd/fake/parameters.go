package main

import (
	"Event-Scheduler/components/common"
	"Event-Scheduler/components/proto"
	"Event-Scheduler/scheduler"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type ContextOptions struct {
	Students      int
	Judges        int
	Rooms         int
	TimeDivisions int

	Events        int
	GroupCapacity int
}

func FakeContext(config ContextOptions) scheduler.ScheduleContext {
	divisions := []int64{}
	for i := 0; i < config.TimeDivisions; i++ {
		divisions = append(divisions, 30)
	}

	students := []*proto.Student{}
	for i := 0; i < config.Students; i++ {
		students = append(students, &proto.Student{
			Email:     gofakeit.Email(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
		})
	}

	judges := []*proto.Judge{}
	for i := 0; i < config.Judges; i++ {
		judges = append(judges, &proto.Judge{
			Email:     gofakeit.Email(),
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
		})
	}

	rooms := []*proto.Room{}
	for i := 0; i < config.Rooms; i++ {
		rooms = append(rooms, &proto.Room{
			Name: fmt.Sprintf(
				"%v%d", strings.ToUpper(gofakeit.Letter()),
				gofakeit.Number(100, 999),
			),
		})
	}

	events := []*proto.Event{}
	for i := 0; i < config.Events; i++ {
		events = append(events, &proto.Event{
			ID: strings.ToUpper(gofakeit.LetterN(3)),
		})
	}

	return scheduler.NewContext(
		&proto.Time{
			Start:     time.Date(2022, time.January, 1, 12, 0, 0, 0, time.Local).Unix(),
			Divisions: divisions,
		},
		&proto.Constraints{
			JudgeStudents: int32(config.TimeDivisions),
			GroupSize:     int32(config.GroupCapacity),
		},
		&proto.Registration{
			Students: students,
			Judges:   judges,
			Rooms:    rooms,
			Events:   events,
		},
	)
}

type RequestOptions struct {
	MinimumEvents int
	MaximumEvents int
	SoloRatio     float64
}

func FakeRequests(c scheduler.ScheduleContext, o RequestOptions) []*proto.StudentRequest {
	requests := []*proto.StudentRequest{}

	eventPool := common.Keys(c.Events)
	studentPool := common.Keys(c.Students)

	for _, studentID := range studentPool {
		amount := o.MinimumEvents + rand.Intn(
			o.MaximumEvents-o.MinimumEvents,
		)
		events := common.SelectRandom(eventPool, amount)

		group := []string{studentID}
		solo := rand.Float64() < o.SoloRatio
		if !solo {
			number := rand.Intn(int(c.GroupSize))
			partners := common.SelectRandom(
				common.Without(studentPool, studentID), number,
			)
			group = append(group, partners...)
		}

		for _, e := range events {
			requests = append(requests, &proto.StudentRequest{
				Event: e,
				Group: group,
			})
		}
	}

	return requests
}
