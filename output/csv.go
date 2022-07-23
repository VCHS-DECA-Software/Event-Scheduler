package output

import (
	"Event-Scheduler/scheduler"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"
)

func CSV(f io.Writer, output scheduler.Output) error {
	writer := csv.NewWriter(f)

	rooms := []string{""}
	for _, h := range output.Housings {
		rooms = append(rooms, h.Room.Name)
		if len(h.Judges) > 0 {
			rooms = append(rooms, make([]string, len(h.Judges)-1)...)
		}
	}
	err := writer.Write(rooms)
	if err != nil {
		return err
	}

	judges := []string{"Timeslots"}
	for _, h := range output.Housings {
		for _, j := range h.Judges {
			accepted := j.Judge.Judgeable
			if len(accepted) == 0 {
				for _, e := range output.Context.Events {
					accepted = append(accepted, e.Id)
				}
			}
			judges = append(
				judges, fmt.Sprintf(
					"%d - %v %v %v",
					j.Judge.Number,
					j.Judge.Firstname,
					j.Judge.Lastname,
					accepted,
				),
			)
		}
	}
	err = writer.Write(judges)
	if err != nil {
		return err
	}

	start := time.Unix(output.Context.Start, 0)
	for i := 0; i < len(output.Context.Divisions); i++ {
		end := start.Add(time.Minute * time.Duration(output.Context.Divisions[i]))
		row := []string{
			fmt.Sprintf(
				"%v - %v",
				start.Format(time.Kitchen),
				end.Format(time.Kitchen),
			),
		}
		start = end

		for _, h := range output.Housings {
			for _, j := range h.Judges {
				if i >= len(j.Assignments) {
					row = append(row, "")
					continue
				}
				names := []string{}
				for _, s := range j.Assignments[i].Group {
					names = append(names, fmt.Sprintf(
						"%v %v",
						s.Firstname,
						s.Lastname,
					))
				}
				if j.Assignments[i].Event == nil {
					row = append(row, "")
					continue
				}
				row = append(row, fmt.Sprintf(
					"%v - %v",
					strings.Join(names, ", "),
					j.Assignments[i].Event.Id,
				))
			}
		}

		err = writer.Write(row)
		if err != nil {
			return err
		}
	}

	writer.Flush()
	return nil
}
