package scheduler

import (
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
			Assignments: []Assignment{},
		})
	}

	leftover := []Assignment{}
	for _, a := range assignments {
		/* this section does the following
		1. appends assignment to the judge if it already contains another
			assignment with the same event
		2. if the judge has nothing assigned it will append the assignment
		(the hypothetical case where an empty-handed judge comes before a
		judge with assignments of the same event cannot happen) */
		assigned := false
		for _, j := range judges {
			if len(j.Assignments) > 0 &&
				len(j.Assignments) < int(c.JudgeStudents) {
				if j.Assignments[len(j.Assignments)-1].Event.ID == a.Event.ID {
					j.Assignments = append(j.Assignments, a)
					assigned = true
					break
				}
			} else if len(j.Assignments) == 0 {
				j.Assignments = append(j.Assignments, a)
				assigned = true
				break
			}
		}

		if !assigned {
			leftover = append(leftover, a)
		}
	}

	/* account for leftover assignments, this will attempt to assign the
	leftovers to any judges with space left, if there are still leftovers
	then report a warning */
	reassigned := 0
	for _, a := range leftover {
		for _, j := range judges {
			if len(j.Assignments) < int(c.JudgeStudents) {
				j.Assignments = append(j.Assignments, a)
				reassigned++
				break
			}
		}
	}
	if reassigned < len(leftover) {
		Warn(fmt.Sprintf(
			"there are %v leftover student requests that exceed "+
				"maximum judge capacity", len(leftover)-reassigned,
		))
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
