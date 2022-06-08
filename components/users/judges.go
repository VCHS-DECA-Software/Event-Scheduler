package users

import (
	"errors"
	"main/components/dbmanager"
	"main/components/events"
)

func GetJudge(id string) (Judge, error) {
	var judge Judge
	err := dbmanager.Query("ID", id, &judge)
	judge.Username = ""
	judge.Password = ""

	eventNames := make([]string, 0)
	for _, eventID := range judge.Events {
		var event events.Event
		dbmanager.Query("ID", eventID, &event)
		eventNames = append(eventNames, event.Name)
	}
	judge.Events = eventNames

	studentNames := make([]string, 0)
	for _, assignedStudentID := range judge.AssignedTeams {
		var student Student
		dbmanager.Query("ID", assignedStudentID, &student)
		student.Password = ""
		studentNames = append(studentNames, student.Name)
	}
	judge.AssignedTeams = studentNames

	return judge, err
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
	judge.Events = append(judge.Events, eventID)
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

	for i, id := range judge.Events {
		if id == eventID {
			judge.Events = append(judge.Events[:i], judge.Events[i+1:]...)
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
	for _, id := range judge.Events {
		if id == eventID {
			return true
		}
	}
	return false
}
