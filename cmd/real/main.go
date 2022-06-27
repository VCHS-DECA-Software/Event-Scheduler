package main

import (
	"Event-Scheduler/components/proto"
	"Event-Scheduler/output"
	"Event-Scheduler/scheduler"
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

var parenthesis = regexp.MustCompile(`\(.+\)`)

func splitName(combined string) (string, string) {
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

func findStudent(pool []*proto.Student, first, last string) *proto.Student {
	var mostLikely *proto.Student
	minFirst := 9999
	minLast := 9999
	for _, s := range pool {
		firstnameDifference := fuzzy.LevenshteinDistance(s.Firstname, first)
		lastnameDifference := fuzzy.LevenshteinDistance(s.Lastname, last)
		if firstnameDifference < minFirst && lastnameDifference < minLast {
			minFirst = firstnameDifference
			minLast = lastnameDifference
			mostLikely = s
		}
	}
	return mostLikely
}

func main() {
	f, err := os.Create("output.log")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)

	f, err = os.Open("real_data.csv")
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(f)

	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	lines = lines[1:]

	students := []*proto.Student{}
	for _, l := range lines {
		email := l[0]
		lastname, firstname := splitName(l[1])

		for _, s := range students {
			if s.Email == email {
				log.Println("duplicate student", firstname, lastname)
				continue
			}
		}

		students = append(students, &proto.Student{
			Email:     email,
			Firstname: firstname,
			Lastname:  lastname,
		})
	}

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

		log.Println("-----------------------", students[i].Firstname, students[i].Lastname)
		group := []string{students[i].Email}
		for _, name := range partners {
			first, last := splitName(name)
			s := findStudent(students, first, last)
			if s == nil {
				log.Printf("could not find student of name %v", name)
				continue
			}
			log.Println(first, last, s.Email)
			group = append(group, s.Email)
		}

		requests = append(requests, &proto.StudentRequest{
			Event: event,
			Group: group,
		})
	}

	log.Println("================================== Students")
	for _, s := range students {
		log.Println(s.Firstname, s.Lastname, s.Email)
	}
	log.Println("================================== Requests")
	for _, r := range requests {
		log.Println(r.Event, r.Group)
	}

	c := scheduler.NewContext(
		&proto.Time{
			Start:     time.Date(2022, time.January, 1, 12, 0, 0, 0, time.Local).Unix(),
			Divisions: []int64{30, 30, 30, 30, 30, 30},
		},
		&proto.Constraints{
			GroupSize: 3,
		},
		&proto.Registration{
			Students: students,
			Judges:   judges,
			Rooms: []*proto.Room{
				{Name: "C410", JudgeCapacity: 3},
				{Name: "C415", JudgeCapacity: 3},
				{Name: "C307", JudgeCapacity: 5},
				{Name: "C309", JudgeCapacity: 5},
				{Name: "Conservatory Hall", JudgeCapacity: 9},
			},
			Events: events,
		},
	)

	log.Println("==================================")
	log.Printf(
		"[STAT] scheduling with %v students, %v requests, and %v judges\n",
		len(students), len(requests), len(judges),
	)

	o := scheduler.Schedule(c, requests)

	f, err = os.Create("output.csv")
	if err != nil {
		panic(err)
	}

	err = output.CSV(f, o)
	if err != nil {
		panic(err)
	}
}
