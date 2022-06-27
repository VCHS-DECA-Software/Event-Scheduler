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

	Events         int
	GroupCapacity  int
	MaxJudgeTalent int
}

func FakeContext(config ContextOptions) scheduler.ScheduleContext {
	divisions := []int64{}
	for i := 0; i < config.TimeDivisions; i++ {
		divisions = append(divisions, 30)
	}

	events := []*proto.Event{}
	for i := 0; i < config.Events; i++ {
		events = append(events, &proto.Event{
			Id: strings.ToUpper(gofakeit.LetterN(3)),
		})
	}

	students := []*proto.Student{}
	for i := 0; i < config.Students; i++ {
		students = append(students, &proto.Student{
			Email:     gofakeit.Email(),
			Firstname: gofakeit.FirstName(),
			Lastname:  gofakeit.LastName(),
		})
	}

	judges := []*proto.Judge{}
	for i := 0; i < config.Judges; i++ {
		number := rand.Intn(int(config.MaxJudgeTalent))
		if number == 0 {
			number = 1
		}

		judgeable := []string{}
		for _, e := range common.SelectRandom(events, int(number)) {
			judgeable = append(judgeable, e.Id)
		}

		judges = append(judges, &proto.Judge{
			Email:     gofakeit.Email(),
			Firstname: gofakeit.FirstName(),
			Lastname:  gofakeit.LastName(),
			Judgeable: judgeable,
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

	return scheduler.NewContext(
		&proto.Time{
			Start:     time.Date(2022, time.January, 1, 12, 0, 0, 0, time.Local).Unix(),
			Divisions: divisions,
		},
		&proto.Constraints{
			GroupSize: int32(config.GroupCapacity),
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

		for _, e := range events {
			group := []string{studentID}
			solo := rand.Float64() < o.SoloRatio
			if !solo {
				number := rand.Intn(int(c.GroupSize))
				partners := common.SelectRandom(
					common.Without(studentPool, studentID), number,
				)
				group = append(group, partners...)
			}
			requests = append(requests, &proto.StudentRequest{
				Event: e,
				Group: group,
			})
		}
	}

	return requests
}
