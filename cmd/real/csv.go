package main

import (
	"Event-Scheduler/components/proto"
	"Event-Scheduler/scheduler"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func FindStudent(pool []*proto.Student, first, last string) (*proto.Student, bool) {
	first = strings.ToLower(first)
	last = strings.ToLower(last)
	for _, s := range pool {
		if first == strings.ToLower(s.Firstname) && last == strings.ToLower(s.Lastname) {
			return s, false
		}
	}
	return &proto.Student{
		Email:     fmt.Sprintf("%v.%v@warriorlife.net", first, last),
		Firstname: first,
		Lastname:  last,
	}, true
}

func SplitName(combined string) (string, string) {
	without := strings.ReplaceAll(string(
		parenthesis.ReplaceAll([]byte(combined), []byte("")),
	), ",", "")
	parts := strings.Split(without, " ")
	normalized := []string{}
	for _, l := range parts {
		if len(l) > 0 {
			normalized = append(normalized, l)
		}
	}
	if len(normalized) == 0 {
		return "", ""
	}
	if len(normalized) == 1 {
		return normalized[0], ""
	}
	return normalized[0], normalized[1]
}

func ParseStudents(lines [][]string) []*proto.Student {
	students := []*proto.Student{}
	for _, l := range lines {
		email := l[0]
		lastname, firstname := SplitName(l[1])

		for _, s := range students {
			if s.Email == email {
				continue
			}
		}

		students = append(students, &proto.Student{
			Email:     email,
			Firstname: firstname,
			Lastname:  lastname,
		})
	}
	return students
}

func ParseRequests(lines [][]string, students []*proto.Student) []*proto.StudentRequest {
	requests := []*proto.StudentRequest{}
	for i, l := range lines {
		event := strings.Split(l[3], " ")[0]

		partners := []string{}
		for _, p := range strings.Split(l[2], ",") {
			trimmed := strings.Trim(p, " ")
			if len(trimmed) > 0 {
				partners = append(partners, trimmed)
			}
		}

		group := []string{students[i].Email}
		for _, name := range partners {
			first, last := SplitName(name)
			s, new := FindStudent(students, first, last)
			if new {
				scheduler.Info(fmt.Sprintf(
					"couldn't find a student by the name \"%v\", "+
						"automatically adding it to the student list...",
					name,
				))
				students = append(students, s)
				continue
			}
			group = append(group, s.Email)
		}

		requests = append(requests, &proto.StudentRequest{
			Event: event,
			Group: group,
		})
	}

	return requests
}

func ParseJudges(rows [][]string) []*proto.Judge {
	judges := []*proto.Judge{}
	for i, row := range rows {
		eventsTrimmed := strings.TrimSpace(row[2])
		events := strings.Split(eventsTrimmed, ",")
		judges = append(judges, &proto.Judge{
			Number:    int32(i),
			Firstname: row[0],
			Lastname:  row[1],
			Judgeable: events,
		})
	}
	return judges
}

func ParseTime(row []string) time.Time {
	startTime, err := time.ParseInLocation("01/02/2006 15:04", row[0], time.Local)
	if err != nil {
		panic(err)
	}
	return startTime
}

func ParseDivisions(rows [][]string) []int64 {
	divisions := []int64{}
	for _, row := range rows {
		if row[0] == "NaN" {
			continue
		}
		slot, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			panic(err)
		}
		divisions = append(divisions, slot)
	}
	return divisions
}

func ParseRooms(rows [][]string) []*proto.Room {
	rooms := []*proto.Room{}
	for _, row := range rows {
		if row[0] == "" {
			continue
		}
		capacity, err := strconv.ParseInt(row[1], 10, 32)
		if err != nil {
			capacity = 0
		}
		rooms = append(rooms, &proto.Room{
			Name:          row[0],
			JudgeCapacity: int32(capacity),
		})
	}
	return rooms
}

func ParseEvents(rows [][]string) []*proto.Event {
	events := []*proto.Event{}
	for _, row := range rows {
		if row[0] == "" {
			continue
		}
		events = append(events, &proto.Event{
			Id: row[0],
		})
	}
	return events
}

func ParseGroupSize(row []string) int64 {
	groupSize, err := strconv.ParseInt(row[0], 10, 32)
	if err != nil {
		panic(err)
	}
	return groupSize
}
