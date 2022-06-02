package users

import (
	"errors"
	"main/components/dbmanager"
	"main/components/encryption"
	"main/components/events"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Student struct {
	ID               string `storm:"id"`
	Username         string `storm:"unique"`
	Name             string
	Password         string
	MyEventIDs      []string
	AssignedJudgeIDs []string
}

func CreateStudent(username, name, password string) error {
	id := uuid.NewV4().String()
	hashedPassword, err := encryption.HashPassword(password)
	if err != nil {
		return err
	}
	err = dbmanager.Save(&Student{ID: id, Username: username, Name: name, Password: hashedPassword})
	return err
}

func GetStudent(id string) (Student, error) {
	var student Student
	err := dbmanager.Query("ID", id, &student)
	return student, err
}

func UpdateStudent(id, username, name, password string) error {
	var student Student
	err := dbmanager.Query("ID", id, &student)
	if err != nil {
		return err
	}
	student.Username = username
	student.Name = name
	hashedPassword, _ := encryption.HashPassword(password)
	student.Password = hashedPassword
	err = dbmanager.Update(&student)
	return err
}

func DeleteStudent(id string) error {
	var student Student
	err := dbmanager.Query("ID", id, &student)
	if err != nil {
		return err
	}
	err = dbmanager.Delete(&student)
	return err
}

func (student Student) JoinEvent(eventID string) error {
	
	if student.IsStudentInEvent(eventID) {
		return errors.New("student is already in event")
	}

	var event events.Event
	err := dbmanager.Query("ID", eventID, &event)
	if err != nil {
		return err
	}

	if len(event.StudentIDs) >= event.MaxStudents {
		return errors.New("event is full")
	}

	if len(student.MyEventIDs) > 3 {
		return errors.New("student has too many events")
	}

	if len(student.MyEventIDs) >= 2 && !strings.Contains(strings.ToLower(event.EventType), "written presentation") {
		return errors.New("student must sign up for a written presentation")
	}

	if student.multipleEventsAtSameTime(event) {
		return errors.New("student is already in another event at the same time")
	}

	student.MyEventIDs = append(student.MyEventIDs, eventID)
	err = dbmanager.Update(&student)
	if err != nil {
		return err
	}
	event.StudentIDs = append(event.StudentIDs, student.ID)
	err = dbmanager.Update(&event)
	return err
}

func (student Student) LeaveEvent(eventID string) error {
	var event events.Event
	err := dbmanager.Query("ID", eventID, &event)
	if err != nil {
		return err
	}

	for i, id := range student.MyEventIDs {
		if id == eventID {
			student.MyEventIDs = append(student.MyEventIDs[:i], student.MyEventIDs[i+1:]...)
			break
		}
	}

	err = dbmanager.Update(&student)
	if err != nil {
		return err
	}

	for i, id := range event.StudentIDs {
		if id == student.ID {
			event.StudentIDs = append(event.StudentIDs[:i], event.StudentIDs[i+1:]...)
			break
		}
	}
	err = dbmanager.Update(&event)
	return err
}

func (student Student) IsStudentInEvent(eventID string) bool {
	for _, id := range student.MyEventIDs {
		if id == eventID {
			return true
		}
	}
	return false
}

func (student Student) multipleEventsAtSameTime(currEvent events.Event) bool {
	for _, eventID := range student.MyEventIDs {
		var event events.Event
		dbmanager.Query("ID", eventID, &event)

		startTime, _ := time.Parse(time.RFC3339, event.StartTime)
		endTime, _ := time.Parse(time.RFC3339, event.EndTime)

		startTimeCurrEvent, _ := time.Parse(time.RFC3339, currEvent.StartTime)
		endTimeCurrEvent, _ := time.Parse(time.RFC3339, currEvent.EndTime)

		if (startTimeCurrEvent.After(startTime) && endTimeCurrEvent.Before(endTime)) || (startTimeCurrEvent.Before(startTime) && endTimeCurrEvent.Before(endTime)) || (startTimeCurrEvent.After(startTime) && endTimeCurrEvent.After(endTime)) || (startTimeCurrEvent.Equal(startTime) || endTimeCurrEvent.Equal(endTime)) {
			return true
		}
	}
	return false
}
