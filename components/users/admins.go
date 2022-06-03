package users

import (
	"main/components/dbmanager"
	"main/components/encryption"
	"math"

	"main/components/events"

	uuid "github.com/satori/go.uuid"
)

type Admin struct {
	ID       string `storm:"id"`
	Username string `storm:"unique"`
	Name     string
	Password string
}

func CreateAdmin(username, name, password string) error {
	id := uuid.NewV4().String()
	hashedPassword, err := encryption.HashPassword(password)
	if err != nil {
		return err
	}
	err = dbmanager.Save(&Admin{ID: id, Username: username, Name: name, Password: hashedPassword})
	return err
}

func GetAdmin(id string) (Admin, error) {
	var admin Admin
	err := dbmanager.Query("ID", id, &admin)
	admin.Password = ""
	return admin, err
}

func UpdateAdmin(id, username, name, password string) error {
	var admin Admin
	err := dbmanager.Query("ID", id, &admin)
	if err != nil {
		return err
	}
	admin.Username = username
	admin.Name = name
	hashedPassword, _ := encryption.HashPassword(password)
	admin.Password = hashedPassword
	err = dbmanager.Update(&admin)
	return err
}

func DeleteAdmin(id string) error {
	var admin Admin
	err := dbmanager.Query("ID", id, &admin)
	if err != nil {
		return err
	}
	err = dbmanager.Delete(&admin)
	return err
}

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
			if len(judge.AssignedStudentIDs) <= 6 {
				judge.AssignedStudentIDs = append(judge.AssignedStudentIDs, team.ID)
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
