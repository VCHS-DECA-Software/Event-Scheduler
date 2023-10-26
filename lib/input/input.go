package input

import (
	"event-scheduler/lib/db/model"
	"event-scheduler/lib/db/table"
	"event-scheduler/lib/utils"
	"fmt"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/go-jet/jet/v2/qrm"
)

type NamePair struct {
	First string
	Last  string
}

func (p NamePair) FormatName() string {
	return fmt.Sprintf("%s %s", p.First, p.Last)
}

func (p NamePair) FormatEmail() string {
	return fmt.Sprintf("%s.%s@warriorlife.net", p.First, p.Last)
}

// takes in a name in the format "last, first" and converts it to "first last"
// it will also convert the name to lowercase
func ParseNamePair(name string) NamePair {
	name = strings.ToLower(name)
	byComma := utils.SplitAndTrim(name, ",")
	if len(byComma) == 2 {
		return NamePair{
			First: byComma[1],
			Last:  byComma[0],
		}
	}
	bySpace := utils.SplitAndTrim(name, " ")
	if len(bySpace) == 2 {
		return NamePair{
			First: bySpace[0],
			Last:  bySpace[1],
		}
	}
	if len(bySpace) == 1 {
		return NamePair{
			First: bySpace[0],
		}
	}
	return NamePair{}
}

const (
	EMPTY_NUMBER = "NaN"
	EMPTY_STRING = ""
)

func getFilteredCol(df dataframe.DataFrame, colname string, not string) series.Series {
	return df.Filter(dataframe.F{
		Colname:    colname,
		Comparator: series.Neq,
		Comparando: not,
	}).Col(colname)
}

func LoadInputs(db qrm.DB, df dataframe.DataFrame) error {
	// constraints
	startTimeStr := df.Col("start_time").Elem(0).String()
	startTime, err := time.Parse(time.Kitchen, startTimeStr)
	if err != nil {
		return err
	}
	maxGroupSize, err := df.Col("group_size").Elem(0).Int()
	if err != nil {
		return err
	}
	examLength, err := df.Col("exam_length").Elem(0).Int()
	if err != nil {
		return err
	}
	config := model.Config{
		ID:           0,
		TimeStart:    startTime.UTC().Format(time.TimeOnly),
		MaxGroupSize: int32(maxGroupSize),
		ExamLength:   int32(examLength),
	}
	_, err = table.Config.INSERT(table.Config.AllColumns).
		MODEL(config).Exec(db)
	if err != nil {
		return err
	}

	timeSlotCount := 0
	for _, v := range df.Col("time_slot").Records() {
		if v == EMPTY_NUMBER {
			break
		}
		timeSlotCount++
	}
	timeSlots, err := df.Col("time_slot").Slice(0, timeSlotCount).Int()
	if err != nil {
		return err
	}
	timeSlotModels := make([]model.TimeSlot, len(timeSlots))
	for i, slot := range timeSlots {
		timeSlotModels[i] = model.TimeSlot{
			ID:       int32(i),
			ConfigID: config.ID,
			Duration: int32(slot),
		}
	}

	_, err = table.TimeSlot.INSERT(table.TimeSlot.AllColumns).
		MODELS(timeSlotModels).Exec(db)
	if err != nil {
		return err
	}

	// rooms
	roomCount := 0
	for _, r := range df.Col("room").Records() {
		if r == EMPTY_STRING {
			break
		}
		roomCount++
	}
	rooms := df.Col("room").Slice(0, roomCount).Records()
	capacities, err := df.Col("room_capacity").Slice(0, len(rooms)).Int()
	if err != nil {
		return err
	}
	roomModels := make([]model.Room, len(rooms))
	for i := 0; i < len(rooms); i++ {
		roomModels[i] = model.Room{
			ID:            rooms[i],
			JudgeCapacity: int32(capacities[i]),
		}
	}
	_, err = table.Room.INSERT(table.Room.AllColumns).
		MODELS(roomModels).Exec(db)
	if err != nil {
		return err
	}

	// events
	eventCount := 0
	for _, e := range df.Col("event").Records() {
		if e == EMPTY_STRING {
			break
		}
		eventCount++
	}
	events := df.Col("event").Slice(0, eventCount).Records()
	eventModels := make([]model.Event, len(events))
	for i, e := range events {
		eventModels[i] = model.Event{
			ID: e,
		}
	}
	_, err = table.Event.INSERT(table.Event.AllColumns).
		MODELS(eventModels).Exec(db)
	if err != nil {
		return err
	}

	// judges
	judgeCount := 0
	for _, name := range df.Col("judge_name").Records() {
		if name == EMPTY_STRING {
			break
		}
		judgeCount++
	}

	judgeNames := df.Col("judge_name").Slice(0, judgeCount).Records()
	judgeCapabilities := df.Col("judge_capability").Slice(0, judgeCount).Records()

	judgeModels := make([]model.Judge, len(judgeNames))
	judgeCapabilityModels := []model.JudgeEventRel{}
	for i := 0; i < len(judgeNames); i++ {
		judgeModels[i] = model.Judge{
			ID:   int32(i),
			Name: judgeNames[i],
		}
		for _, eventId := range utils.SplitAndTrim(judgeCapabilities[i], ",") {
			judgeCapabilityModels = append(judgeCapabilityModels, model.JudgeEventRel{
				JudgeID: int32(i),
				EventID: eventId,
			})
		}
	}
	_, err = table.Judge.INSERT(table.Judge.AllColumns).
		MODELS(judgeModels).Exec(db)
	if err != nil {
		return err
	}
	_, err = table.JudgeEventRel.INSERT(table.JudgeEventRel.AllColumns).
		MODELS(judgeCapabilityModels).Exec(db)
	if err != nil {
		return err
	}

	// students
	studentEmails := getFilteredCol(df, "student_email", "").Records()
	studentNames := getFilteredCol(df, "student_name", "").Records()
	groupPartners := df.Col("group_partners").Records()
	groupEvents := getFilteredCol(df, "group_event", "").Records()

	studentModels := make([]model.Student, len(studentEmails))
	for i := 0; i < len(studentEmails); i++ {
		namePair := ParseNamePair(studentNames[i])
		studentModels[i] = model.Student{
			Email: studentEmails[i],
			Name:  namePair.FormatName(),
		}
	}
	_, err = table.Student.
		INSERT(table.Student.AllColumns).
		MODELS(studentModels).
		ON_CONFLICT().
		DO_NOTHING().
		Exec(db)
	if err != nil {
		return err
	}

	groupModels := make([]model.StudentGroup, len(studentEmails))
	partnerModels := []model.StudentStudentGroupRel{}
	for i := 0; i < len(studentEmails); i++ {
		event := groupEvents[i]
		groupModels[i] = model.StudentGroup{
			ID:      int32(i),
			EventID: event,
		}

		if groupPartners[i] == "" {
			continue
		}
		partners := utils.SplitAndTrim(groupPartners[i], ",")
		partnerModels = append(partnerModels, model.StudentStudentGroupRel{
			StudentGroupID: int32(i),
			StudentEmail:   studentEmails[i],
		})

		for _, partner := range partners {
			namePair := ParseNamePair(partner)
			partnerEmail := namePair.FormatEmail()
			_, err = table.Student.INSERT(table.Student.AllColumns).
				MODEL(model.Student{
					Email: partnerEmail,
					Name:  namePair.FormatName(),
				}).ON_CONFLICT().DO_NOTHING().Exec(db)
			if err != nil {
				return err
			}
			partnerModels = append(partnerModels, model.StudentStudentGroupRel{
				StudentGroupID: int32(i),
				StudentEmail:   partnerEmail,
			})
		}
	}
	_, err = table.StudentGroup.
		INSERT(table.StudentGroup.AllColumns).
		MODELS(groupModels).
		Exec(db)
	if err != nil {
		return err
	}
	_, err = table.StudentStudentGroupRel.
		INSERT(table.StudentStudentGroupRel.AllColumns).
		MODELS(partnerModels).
		Exec(db)
	if err != nil {
		return err
	}

	return nil
}
