package scheduler

import (
	"Event-Scheduler/components/common"
	"Event-Scheduler/components/proto"
	"fmt"
	"log"
	"math"
)

func Warn(message string) {
	log.Printf("[WARN] %v\n", message)
}

type ScheduleContext struct {
	*proto.Time
	*proto.Constraints

	Students map[string]*proto.Student
	Judges   map[string]*proto.Judge
	Events   map[string]*proto.Event
	Rooms    []*proto.Room
}

func NewContext(
	t *proto.Time,
	c *proto.Constraints,
	r *proto.Registration,
) ScheduleContext {
	context := ScheduleContext{
		Time:        t,
		Constraints: c,
		Students:    map[string]*proto.Student{},
		Judges:      map[string]*proto.Judge{},
		Events:      map[string]*proto.Event{},
		Rooms:       r.Rooms,
	}

	for _, s := range r.Students {
		context.Students[s.Email] = s
	}
	for _, j := range r.Judges {
		context.Judges[j.Email] = j
	}
	for _, e := range r.Events {
		context.Events[e.ID] = e
	}

	return context
}

func Schedule(c ScheduleContext, requests []*proto.StudentRequest) Output {
	assignments := []Assignment{}
	for _, r := range requests {
		group := []*proto.Student{}
		for _, student := range r.Group {
			group = append(group, c.Students[student])
		}
		assignments = append(assignments, Assignment{
			Group: group,
			Event: c.Events[r.Event],
		})
	}

	//initialize judge structs
	judges := []*Judgement{}
	for _, j := range c.Judges {
		judges = append(judges, &Judgement{
			Judge:       j,
			Assignments: make([]Assignment, len(c.Divisions)),
		})
	}

	leftover := []Assignment{}
assignments:
	for _, a := range assignments {
		//see "algorithm" in docs/scheduling.md
		occupied := map[int]bool{}
		for _, j := range judges {
			for i := 0; i < len(c.Divisions); i++ {
				if common.Intersects(j.Assignments[i].Group, a.Group) {
					occupied[i] = true
				}
			}
		}

		for _, j := range judges {
			for i := 0; i < len(c.Divisions); i++ {
				if occupied[i] {
					continue
				}
				if j.Assignments[i].Event != nil {
					continue
				}
				//checks if there is an (vertically) adjacent
				//assignment with the same event
				if (i > 0 && j.Assignments[i-1].Event != nil &&
					j.Assignments[i-1].Event.ID == a.Event.ID) ||
					(i < len(c.Divisions)-1 && j.Assignments[i+1].Event != nil &&
						j.Assignments[i+1].Event.ID == a.Event.ID) {
					j.Assignments[i] = a
					continue assignments
				}
			}
		}

		for _, j := range judges {
			for i := 0; i < len(c.Divisions); i++ {
				if occupied[i] {
					continue
				}
				if j.Assignments[i].Event == nil {
					j.Assignments[i] = a
					continue assignments
				}
			}
		}

		leftover = append(leftover, a)
	}

	if len(leftover) > 0 {
		Warn(fmt.Sprintf(
			"there are %v leftover student requests that could not "+
				"be assigned without conflicts", len(leftover),
		))
	}

	for i := 0; i < len(c.Divisions); i++ {
		contains := map[string]bool{}
		for _, j := range judges {
			if i >= len(j.Assignments) {
				continue
			}
			for _, s := range j.Assignments[i].Group {
				if !contains[s.Email] {
					contains[s.Email] = true
					continue
				}
				Warn(fmt.Sprintf(
					"there is a conflict involving %v on division %v",
					s.FirstName, i,
				))
			}
		}
	}

	//try and spread out judges evenly throughout the rooms
	roomSize := int(math.Ceil(float64(len(judges)) / float64(len(c.Rooms))))
	housings := []Housing{}
	for i := 0; i < len(c.Rooms); i++ {
		start := i * roomSize
		if start < 0 {
			start = 0
		}
		end := (i + 1) * roomSize
		if end > len(judges) {
			end = len(judges)
		}
		housings = append(housings, Housing{
			Room:   c.Rooms[i],
			Judges: judges[start:end],
		})
	}

	return Output{
		Housings: housings,
		Context:  c,
	}
}
