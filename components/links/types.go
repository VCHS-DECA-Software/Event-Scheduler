package links

import uuid "main/vendor/github.com/satori/go.uuid"

type EventType = int

const (
	ORAL = iota
	ROLEPLAY
	EXAM
)

type Event struct {
	ID string `storm:"id"`

	Judges Link
	Teams  Link

	AdminID string `storm:"index"`
	Name    string `storm:"unique"`

	Description string
	Location    string    `storm:"index"`
	EventType   EventType `storm:"index"`

	StartTime string
	EndTime   string
}

type Team struct {
	ID string `storm:"id"`

	Students Link
	Events   Link

	OwnerID string
	Name    string `storm:"index"`
}

func NewTeam(owner, name string) Team {
	id := uuid.NewV4().String()
	return Team{
		ID:       id,
		Students: NewLink(-1),
		Events:   NewLink(2),
	}
}
