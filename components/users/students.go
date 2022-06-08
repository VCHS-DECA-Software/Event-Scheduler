package users

import (
	"errors"
	"fmt"
	"main/components/dbmanager"
	"main/components/events"
	"main/components/links"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

func NewStudent(name string) Student {
	return Student{
		Name:   name,
		Teams:  links.NewLink(1),
	}
}

func GetStudent(id string) (Student, error) {
	var student Student
	err := dbmanager.Query("ID", id, &student)
	student.Username = ""
	student.Password = ""
	student.Teams = nil

	eventNames := make([]string, 0)
	for _, eventID := range student.Events {
		var event events.Event
		dbmanager.Query("ID", eventID, &event)
		eventNames = append(eventNames, event.Name)
	}
	student.Events = eventNames
	return student, err
}

func (student *Account[Student]) CreateTeam(name string) error {
	student.Update(Account[Student]{
		Data: Student{},
	})
	id := uuid.NewV4().String()
	err := dbmanager.Save(&Team{ID: id, Name: name, TeamOwnerID: student.ID})
	if err != nil {
		return err
	}
	err = student.JoinTeam(id)
	return err
}

func (student Student) GetTeams() ([]Team, error) {
	var teams []Team
	for _, teamID := range student.Teams {
		team, _ := GetTeam(teamID)
		teams = append(teams, team)
	}
	return teams, nil
}

func readTeam(id string) (Team, error) {
	var team Team
	err := dbmanager.Query("ID", id, &team)
	return team, err
}

func GetTeam(id string) (Team, error) {
	var team Team
	err := dbmanager.Query("ID", id, &team)
	studentNames, _ := viewTeamMembers(team.ID)
	team.StudentIDs = studentNames
	teamOwner, _ := readStudent(team.TeamOwnerID)
	team.TeamOwnerID = teamOwner.Name
	return team, err
}

func GetTeamsbyName(name string) ([]Team, error) {
	var teams []Team
	err := dbmanager.GroupQuery("Name", name, &teams)
	for _, team := range teams {
		studentNames, _ := viewTeamMembers(team.ID)
		team.StudentIDs = studentNames
		teamOwner, _ := readStudent(team.TeamOwnerID)
		team.TeamOwnerID = teamOwner.Name
	}
	return teams, err
}

func viewTeamMembers(id string) ([]string, error) {
	var team Team
	err := dbmanager.Query("ID", id, &team)
	if err != nil {
		return nil, err
	}
	for i, studentID := range team.StudentIDs {
		var student Student
		dbmanager.Query("ID", studentID, &student)
		team.StudentIDs[i] = student.Name
	}
	return team.StudentIDs, nil
}

//ViewAllTeams returns all teams in the database (should be used for admins and judges)
func ViewAllTeams() ([]Team, error) {
	var teams []Team
	err := dbmanager.QueryAll(&teams)
	if err != nil {
		return nil, err
	}

	for _, team := range teams {
		studentNames, _ := viewTeamMembers(team.ID)
		team.StudentIDs = studentNames
		teamOwner, _ := readStudent(team.TeamOwnerID)
		team.TeamOwnerID = teamOwner.Name
	}

	return teams, nil
}

func (student Student) UpdateTeam(id, name, teamOwnerID string) error {
	var team Team
	err := dbmanager.Query("ID", id, &team)
	if err != nil {
		return err
	}

	if !(student.ID == team.TeamOwnerID) {
		return errors.New("only team owners can update teams")
	}

	team.Name = name
	if teamOwnerID != "" {
		team.TeamOwnerID = teamOwnerID
	}
	err = dbmanager.Update(&team)
	return err
}

func (student Student) DeleteTeam(id string) error {
	var team Team
	err := dbmanager.Query("ID", id, &team)
	if err != nil {
		return err
	}

	if !(student.ID == team.TeamOwnerID) {
		return errors.New("only team owners can delete teams")
	}

	err = dbmanager.Delete(&team)
	return err
}

func (student Student) JoinTeam(teamID string) error {
	var team Team
	err := dbmanager.Query("ID", teamID, &team)
	if err != nil {
		return err
	}
	team.StudentIDs = append(team.StudentIDs, student.ID)
	err = dbmanager.Update(&team)
	return err
}

func (student Student) LeaveTeam(teamID string) error {
	var team Team
	err := dbmanager.Query("ID", teamID, &team)
	if err != nil {
		return err
	}
	for i, studentID := range team.StudentIDs {
		if studentID == student.ID {
			team.StudentIDs = append(team.StudentIDs[:i], team.StudentIDs[i+1:]...)
			break
		}
	}
	err = dbmanager.Update(&team)
	return err
}

func (student Student) JoinEvent(teamID, eventID string) error {

	var team Team
	err := dbmanager.Query("ID", teamID, &team)
	if err != nil {
		return err
	}

	if !(student.ID == team.TeamOwnerID) {
		return errors.New("only team owners can sign up teams for events")
	}

	if team.isTeamInEvent(eventID) {
		return errors.New("team is already in event")
	}

	var event events.Event
	err = dbmanager.Query("ID", eventID, &event)
	if err != nil {
		return err
	}

	count, err := countStudentsForEvent(eventID)
	if err != nil {
		return err
	}

	if count >= event.MaxStudents {
		return errors.New("event is full")
	}

	for _, studentID := range team.StudentIDs {
		var student Student
		dbmanager.Query("ID", studentID, &student)

		if !student.checkEventTypes(event) {
			return fmt.Errorf("student %s cannot register for this event", student.Name)
		}

		if student.multipleEventsAtSameTime(event) {
			return fmt.Errorf("student %s has already registered for an event at the same time", student.Name)
		}

		student.Events = append(student.Events, eventID)
		err = dbmanager.Update(&student)
		if err != nil {
			return err
		}
	}

	event.TeamIDs = append(event.TeamIDs, team.ID)
	err = dbmanager.Update(&event)
	return err
}

func (student Student) LeaveEvent(teamID, eventID string) error {

	var team Team
	err := dbmanager.Query("ID", teamID, &team)
	if err != nil {
		return err
	}

	if !(student.ID == team.TeamOwnerID) {
		return errors.New("only team owners can leave events for teams")
	}

	var event events.Event
	err = dbmanager.Query("ID", eventID, &event)
	if err != nil {
		return err
	}

	for i, id := range student.Events {
		if id == eventID {
			student.Events = append(student.Events[:i], student.Events[i+1:]...)
			break
		}
	}

	err = dbmanager.Update(&student)
	if err != nil {
		return err
	}

	for i, id := range event.TeamIDs {
		if id == team.ID {
			event.TeamIDs = append(event.TeamIDs[:i], event.TeamIDs[i+1:]...)
			break
		}
	}
	err = dbmanager.Update(&event)
	return err
}

func (team Team) isTeamInEvent(eventID string) bool {
	event, err := events.GetEvent(eventID)
	if err != nil {
		return false
	}
	for _, id := range event.TeamIDs {
		if id == team.ID {
			return true
		}
	}
	return false
}

func (student Student) multipleEventsAtSameTime(currEvent events.Event) bool {
	for _, eventID := range student.Events {
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

func countStudentsForEvent(id string) (int, error) {
	var event events.Event
	err := dbmanager.Query("ID", id, &event)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, teamID := range event.TeamIDs {
		var team Team
		err := dbmanager.Query("ID", teamID, &team)
		if err != nil {
			return 0, err
		}
		count += len(team.StudentIDs)
	}
	return count, nil
}

func (student Student) checkEventTypes(currEvent events.Event) bool {

	if len(student.Events) == 0 {
		return true
	}

	if len(student.Events) >= 2 {
		return false
	}

	var firstEvent events.Event
	dbmanager.Query("ID", student.Events[0], &firstEvent)

	if (len(student.Events) == 1) && (strings.Contains(strings.ToLower(firstEvent.EventType), "written presentation")) && (strings.Contains(strings.ToLower(currEvent.EventType), "oral presentation")) {
		return true
	}

	if (len(student.Events) == 1) && (strings.Contains(strings.ToLower(firstEvent.EventType), "oral presentation")) && (strings.Contains(strings.ToLower(currEvent.EventType), "written presentation")) {
		return true
	}

	return false

}
