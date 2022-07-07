package main

import (
	"Event-Scheduler/components/proto"
	"Event-Scheduler/output"
	"Event-Scheduler/scheduler"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
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

func findStudent(pool []*proto.Student, first, last string) (*proto.Student, bool) {
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

func main() {
	f, err := os.Create("output.log")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	log.SetOutput(f)

	studentRegistrationFilePtr := flag.String("student", "", "student registration file")
	judgeRegistrationFilePtr := flag.String("judge", "", "judge registration file")
	conferenceFilePtr := flag.String("conference", "", "conference details file")

	flag.Parse()

	f, err = os.Open(*studentRegistrationFilePtr)
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

		group := []string{students[i].Email}
		for _, name := range partners {
			first, last := splitName(name)
			s, new := findStudent(students, first, last)
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

	log.Println("================================== Students")
	for _, s := range students {
		log.Println(s.Firstname, s.Lastname, s.Email)
	}
	log.Println("================================== Requests")
	for _, r := range requests {
		log.Println(r.Event, r.Group)
	}

	jFile, err := os.Open(*judgeRegistrationFilePtr)
	if err != nil {
		panic(err)
	}

	cFile, err := os.Open(*conferenceFilePtr)
	if err != nil {
		panic(err)
	}

	var jReader io.Reader = jFile
	var cReader io.Reader = cFile

	judgeDf := dataframe.ReadCSV(jReader)
	conferenceDf := dataframe.ReadCSV(cReader)

	judges := []*proto.Judge{}

	for _, row := range judgeDf.Records()[1:] {
		eventsTrimmed := strings.TrimSpace(row[2])
		events := strings.Split(eventsTrimmed, ",")
		judges = append(judges, &proto.Judge{
			Firstname: row[0],
			Lastname:  row[1],
			Judgeable: events,
		})
	}

	conferenceStartTime, err := time.Parse("01/02/2006 15:04", conferenceDf.Select([]string{"Start Time"}).Records()[1][0])

	if err != nil {
		panic(err)
	}

	divisions := []int64{}
	for _, row := range conferenceDf.Select([]string{"Time Slot"}).Records()[1:] {
		slot, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			panic(err)
		}

		divisions = append(divisions, slot)
	}

	rooms := []*proto.Room{}
	for _, row := range conferenceDf.Select([]string{"Room", "Judge Capacity"}).Records()[1:] {
		capacity, err := strconv.ParseInt(row[1], 10, 32)
		if err != nil {
			capacity = 0
		}
		rooms = append(rooms, &proto.Room{
			Name:          row[0],
			JudgeCapacity: int32(capacity),
		})
	}

	events := []*proto.Event{}
	for _, row := range conferenceDf.Select([]string{"Event"}).Records()[1:] {
		events = append(events, &proto.Event{
			Id: row[0],
		})
	}

	groupsizeStr := conferenceDf.Select([]string{"Group Size"}).Records()[1][0]
	groupsize, err := strconv.ParseInt(groupsizeStr, 10, 32)
	if err != nil {
		panic(err)
	}

	c := scheduler.NewContext(
		&proto.Time{
			Start:    conferenceStartTime.Unix(),
			Divisions: divisions,
		},
		&proto.Constraints{
			GroupSize: int32(groupsize),
		},
		&proto.Registration{
			Students: students,
			Judges:   judges,
			Rooms:    rooms,
			Events:   events,
		},
	)

	log.Println("==================================")
	log.Printf(
		"[INFO] scheduling with %v students, %v requests, and %v judges\n",
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
