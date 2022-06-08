package users

import (
	"main/components/dbmanager"
	"math"

	"main/components/events"

	uuid "github.com/satori/go.uuid"
)

func (admin Admin) CreateEvent(name, description, eventType, startTime, endTime, location string, maxStudents int) error {
	id := uuid.NewV4().String()

	if maxStudents == 0 {
		maxStudents = math.MaxInt64
	}

	err := dbmanager.Save(&events.Event{ID: id, Name: name, AdminID: admin.ID, Description: description, EventType: eventType, StartTime: startTime, EndTime: endTime, Location: location, MaxStudents: maxStudents})
	return err
}

func ViewAllEvents() error {
	events, err := events.GetAllEvents()
	if err != nil {
		return err
	}
	for _, event := range events {
		var admin Admin
		dbmanager.Query("ID", event.AdminID, &admin)
		event.AdminID = admin.Name

		for i, judgeID := range event.JudgeIDs {
			var judge Judge
			dbmanager.Query("ID", judgeID, &judge)
			event.JudgeIDs[i] = judge.Name
		}

		for i, teamID := range event.TeamIDs {
			var team Team
			dbmanager.Query("ID", teamID, &team)
			event.TeamIDs[i] = team.Name
		}
	}
	return nil
}

func (admin Admin) UpdateEvent(id, name, description, eventType, startTime, endTime, location string, maxStudents int) error {
	var event events.Event
	err := dbmanager.Query("ID", id, &event)
	if err != nil {
		return err
	}
	event.Name = name
	event.Description = description
	event.EventType = eventType
	event.StartTime = startTime
	event.EndTime = endTime
	event.Location = location

	if maxStudents == 0 {
		maxStudents = math.MaxInt64
	}

	err = dbmanager.Update(&event)
	return err
}

func (admin Admin) DeleteEvent(id string) error {
	var event events.Event
	err := dbmanager.Query("ID", id, &event)
	if err != nil {
		return err
	}
	err = dbmanager.Delete(&event)
	return err
}

func (admin Admin) AssignStudentsAndJudges(id string) error {
	var event events.Event
	err := dbmanager.Query("ID", id, &event)

	teams := make([]Team, 0)
	for _, teamID := range event.TeamIDs {
		team, _ := readTeam(teamID)
		teams = append(teams, team)
	}

	judges := make([]Judge, 0)
	for _, judgeID := range event.JudgeIDs {
		judge, _ := readJudge(judgeID)
		judges = append(judges, judge)
	}

	for _, judge := range judges {
		for _, team := range teams {
			if len(judge.AssignedTeams) <= 6 {
				judge.AssignedTeams = append(judge.AssignedTeams, team.ID)
				team.AssignedJudgeIDs = append(team.AssignedJudgeIDs, judge.ID)
			}
		}
	}

	return err
}

func (admin Admin) AssignForAllEvents() error {
	events, err := events.GetAllEvents()

	for _, event := range events {
		admin.AssignStudentsAndJudges(event.ID)
	}

	return err
}
