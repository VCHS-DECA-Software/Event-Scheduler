package users

import (
	"errors"
	"main/components/dbmanager"
	"main/components/encryption"
	"main/components/events"

	uuid "github.com/satori/go.uuid"
)

type Judge struct {
	ID                 string `storm:"id"`
	Username           string `storm:"unique"`
	Name               string
	Password           string
	MyEventIDs         []string
	AssignedStudentIDs []string
}

func CreateJudge(username, name, password string) error {
	id := uuid.NewV4().String()
	hashedPassword, err := encryption.HashPassword(password)
	if err != nil {
		return err
	}
	err = dbmanager.Save(&Judge{ID: id, Username: username, Name: name, Password: hashedPassword})
	return err
}

func GetJudge(id string) (Judge, error) {
	var judge Judge
	err := dbmanager.Query("ID", id, &judge)
	judge.Username = ""
	judge.Password = ""

	eventNames := make([]string, 0)
	for _, eventID := range judge.MyEventIDs {
		var event events.Event
		dbmanager.Query("ID", eventID, &event)
		eventNames = append(eventNames, event.Name)
	}
	judge.MyEventIDs = eventNames

	studentNames := make([]string, 0)
	for _, assignedStudentID := range judge.AssignedStudentIDs {
		var student Student
		dbmanager.Query("ID", assignedStudentID, &student)
		student.Password = ""
		studentNames = append(studentNames, student.Name)
	}
	judge.AssignedStudentIDs = studentNames

	return judge, err
}

func AuthenticateJudge(username, password string) bool {
	var judge Judge
	err := dbmanager.Query("Username", username, &judge)
	if err != nil {
		return false
	}
	return encryption.CheckPasswordHash(password, judge.Password)
}

func readJudge(id string) (Judge, error) {
	var judge Judge
	err := dbmanager.Query("ID", id, &judge)
	return judge, err
}

func UpdateJudge(id, username, name, password string) error {
	var judge Judge
	err := dbmanager.Query("ID", id, &judge)
	if err != nil {
		return err
	}
	judge.Username = username
	judge.Name = name
	hashedPassword, _ := encryption.HashPassword(password)
	judge.Password = hashedPassword
	err = dbmanager.Update(&judge)
	return err
}

func DeleteJudge(id string) error {
	var judge Judge
	err := dbmanager.Query("ID", id, &judge)
	if err != nil {
		return err
	}
	err = dbmanager.Delete(&judge)
	return err
}

func (judge Judge) JoinEvent(eventID string) error {

	if judge.IsJudgeInEvent(eventID) {
		return errors.New("judge is already in event")
	}

	var event events.Event
	err := dbmanager.Query("ID", eventID, &event)
	if err != nil {
		return err
	}
	judge.MyEventIDs = append(judge.MyEventIDs, eventID)
	err = dbmanager.Update(&judge)
	if err != nil {
		return err
	}
	event.JudgeIDs = append(event.JudgeIDs, judge.ID)
	err = dbmanager.Update(&event)
	return err
}

func (judge Judge) LeaveEvent(eventID string) error {
	var event events.Event
	err := dbmanager.Query("ID", eventID, &event)
	if err != nil {
		return err
	}

	for i, id := range judge.MyEventIDs {
		if id == eventID {
			judge.MyEventIDs = append(judge.MyEventIDs[:i], judge.MyEventIDs[i+1:]...)
			break
		}
	}

	err = dbmanager.Update(&judge)
	if err != nil {
		return err
	}

	for i, id := range event.JudgeIDs {
		if id == judge.ID {
			event.JudgeIDs = append(event.JudgeIDs[:i], event.JudgeIDs[i+1:]...)
			break
		}
	}
	err = dbmanager.Update(&event)
	return err
}

func (judge Judge) IsJudgeInEvent(eventID string) bool {
	for _, id := range judge.MyEventIDs {
		if id == eventID {
			return true
		}
	}
	return false
}
