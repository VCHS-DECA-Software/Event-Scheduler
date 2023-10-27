package scheduler

import (
	"Event-Scheduler/components/common"
	"Event-Scheduler/components/proto"
	"fmt"
	"log"
	"sort"
	"strings"
)

func Info(message string) {
	log.Printf("[INFO] %v\n", message)
}

func Warn(message string) {
	log.Printf("[WARN] %v\n", message)
}

func FormatEvent(event *proto.Event) string {
	return fmt.Sprintf("\"%v\"", event.Id)
}

func FormatGroup(group []*proto.Student) string {
	groupText := []string{}
	for _, s := range group {
		groupText = append(groupText, fmt.Sprintf(
			"%v %v", s.Firstname, s.Lastname,
		))
	}
	return fmt.Sprintf("[%v]", strings.Join(groupText, ", "))
}

type ScheduleContext struct {
	*proto.Time
	*proto.Constraints

	Students map[string]*proto.Student
	Judges   map[int]*proto.Judge
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
		Judges:      map[int]*proto.Judge{},
		Events:      map[string]*proto.Event{},
		Rooms:       r.Rooms,
	}

	for _, s := range r.Students {
		context.Students[s.Email] = s
	}
	for _, j := range r.Judges {
		context.Judges[int(j.Number)] = j
	}
	for _, e := range r.Events {
		context.Events[e.Id] = e
	}

	return context
}

func Schedule(c ScheduleContext, requests []*proto.StudentRequest) Output {
	excludedFromExams := make(map[string]bool)

	assignments := []Assignment{}
assignments:
	for _, r := range requests {
		group := []*proto.Student{}
		for _, student := range r.Group {
			s, ok := c.Students[student]
			if !ok {
				Info(fmt.Sprintf(
					"group's partner (%v) does not exist, skipping...",
					student,
				))
				continue
			}
			group = append(group, s)
		}

		event, ok := c.Events[r.Event]
		if !ok {
			Info(fmt.Sprintf(
				"assignment's target event (%v) does not exist, skipping...\n",
				r.Event,
			))
			continue
		}

		if event.EventType == proto.EventType_WRITTEN {
			for _, student := range r.Group {
				excludedFromExams[student] = true
			}
		}

		for _, a := range assignments {
			if common.UnorderedEqual(a.Group, group) && a.Event == event {
				Info(fmt.Sprintf(
					"duplicate student requests (%v - %v) skipping...",
					FormatEvent(event), FormatGroup(a.Group),
				))
				continue assignments
			}
		}

		assignments = append(assignments, Assignment{
			Group: group,
			Event: event,
		})
	}

	//sort requests from the largest group to the smallest group
	sort.SliceStable(assignments, func(i, j int) bool {
		return len(assignments[i].Group) > len(assignments[j].Group)
	})

	//initialize judge structs
	judges := []*Judgement{}
	for _, j := range c.Judges {
		judges = append(judges, &Judgement{
			Judge:       j,
			Assignments: make([]Assignment, len(c.Divisions)),
		})
	}

	//sort judges from the least flexible to the most flexible
	sort.SliceStable(judges, func(i, j int) bool {
		return len(judges[i].Judge.Judgeable) < len(judges[j].Judge.Judgeable)
	})

	calculateOccupied := func(group []*proto.Student) map[int]bool {
		occupied := map[int]bool{}
		for _, j := range judges {
			for i := 0; i < len(c.Divisions); i++ {
				if common.Intersects(j.Assignments[i].Group, group) {
					occupied[i] = true
				}
			}
		}
		return occupied
	}

	assign := func(occupied map[int]bool, a Assignment, strict bool) bool {
		for _, j := range judges {
			if !common.Intersects([]string{a.Event.Id}, j.Judge.Judgeable) &&
				len(j.Judge.Judgeable) > 0 {
				continue
			}
			for i := 0; i < len(c.Divisions); i++ {
				if occupied[i] {
					continue
				}
				if j.Assignments[i].Event != nil {
					continue
				}

				backToBack := false

				for _, j := range judges {
					if common.HasAdjacent(j.Assignments, i, func(adj Assignment, above bool) bool {
						return common.Intersects(adj.Group, a.Group)
					}) {
						backToBack = true
						break
					}
				}

				if backToBack {
					continue
				}
				if strict {
					//checks if there is an (vertically) adjacent
					//assignment with the same event
					if common.HasAdjacent(j.Assignments, i, func(adj Assignment, above bool) bool {
						return adj.Event != nil && adj.Event.Id == a.Event.Id
					}) {
						j.Assignments[i] = a
						return true
					}
					return false
				}

				j.Assignments[i] = a
				return true
			}
		}
		return false
	}

	leftover := []Assignment{}
	for _, a := range assignments {
		//see "algorithm" in docs/scheduling.md
		occupied := calculateOccupied(a.Group)

		assigned := assign(occupied, a, true)
		if assigned {
			continue
		}
		assigned = assign(occupied, a, false)
		if assigned {
			continue
		}

		leftover = append(leftover, a)
	}

	//this isn't very efficient, but I'm not a competitive programmer
	//	nor am I attempting to pass an interview at google.
	exams := []Exam{}
students:
	for _, s := range c.Students {
		if excludedFromExams[s.Email] {
			continue
		}

		occupied := calculateOccupied([]*proto.Student{s})

		start := 0
		sum := int32(0)
		for i := 0; i < len(c.Divisions); i++ {
			if _, in := occupied[i]; !in {
				sum += c.Divisions[i]
				if sum >= c.Constraints.ExamLength {
					exams = append(exams, Exam{
						Start:   start,
						Student: s,
					})
					continue students
				}
				continue
			}
			sum = 0
			start = i + 1
		}

		Warn(fmt.Sprintf(
			"could not find a suitable exam time for student %v %v",
			s.Firstname, s.Lastname,
		))
	}

	if len(leftover) > 0 {
		Warn(fmt.Sprintf(
			"there are %v leftover student requests that could not "+
				"be assigned without conflicts", len(leftover),
		))
		for _, s := range leftover {
			log.Println(FormatEvent(s.Event), FormatGroup(s.Group))
		}

		numerator := 0
		noJudge := map[string]bool{}
		for _, s := range leftover {
			judgeable := false
			for _, j := range judges {
				if common.Intersects(j.Judge.Judgeable, []string{s.Event.Id}) {
					judgeable = true
					break
				}
			}
			if !judgeable {
				numerator++
				noJudge[s.Event.Id] = true
			}
		}
		Warn(fmt.Sprintf(
			"%v%% of leftover requests are due to having "+
				"no judges able to judge %v",
			(float64(numerator)/float64(len(leftover)))*100,
			common.Keys(noJudge),
		))
	}

	//! DEBUG: conflict checking
	// for i := 0; i < len(c.Divisions); i++ {
	// 	contains := map[string]bool{}
	// 	for _, j := range judges {
	// 		if i >= len(j.Assignments) {
	// 			continue
	// 		}
	// 		for _, s := range j.Assignments[i].Group {
	// 			if !contains[s.Email] {
	// 				contains[s.Email] = true
	// 				continue
	// 			}
	// 			Warn(fmt.Sprintf(
	// 				"there is a conflict involving %v on division %v",
	// 				s.Firstname, i,
	// 			))
	// 		}
	// 	}
	// }

	typedJudgeSet := make(map[proto.EventType][]*Judgement)
	typedRoomSet := make(map[proto.EventType][]*proto.Room)
	for _, judge := range judges {
		judgeEventType := c.Events[judge.Judge.Judgeable[0]].EventType
		typedJudgeSet[judgeEventType] = append(typedJudgeSet[judgeEventType], judge)
	}
	for _, room := range c.Rooms {
		typedRoomSet[room.EventType] = append(typedRoomSet[room.EventType], room)
	}

	//try and spread out judges evenly throughout the rooms
	housings := make(map[proto.EventType][]Housing)
	for eventType, rooms := range typedRoomSet {
		judges := typedJudgeSet[eventType]
		housings[eventType] = make([]Housing, len(rooms))

		roomIndex := 0
		judgeIndex := 0
		filledExplored := 0
		for judgeIndex < len(judges) {
			room := rooms[roomIndex]
			housing := housings[eventType][roomIndex]
			if housing.Room == nil {
				housing.Room = room
			}

			if len(housing.Judges) < int(housing.Room.JudgeCapacity) {
				housing.Judges = append(housing.Judges, judges[judgeIndex])
				housings[eventType][roomIndex] = housing
				judgeIndex++
				filledExplored = 0
			} else {
				filledExplored++
				if filledExplored == len(rooms) {
					Warn(fmt.Sprintf(
						"there is not enough room to house all the judges for the event type \"%s\""+
							"try adjusting 'Judge Capacity', %d judges will be dropped",
						eventType.String(), len(judges)-judgeIndex,
					))
				}
			}

			roomIndex++
			if roomIndex == len(rooms) {
				roomIndex = 0
			}
		}
	}

	return Output{
		Housings: housings,
		Context:  c,
		Exams:    exams,
	}
}
