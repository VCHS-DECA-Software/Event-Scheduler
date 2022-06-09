package groupings

type EventType = int

const (
	ORAL = iota
	ROLEPLAY
	EXAM
)

type Hierarchy[T any] struct {
	Owner string `storm:"index"`
	Name  string `storm:"unique"`
	Data  T
}

type Event struct {
	Description string
	Location    string    `storm:"index"`
	EventType   EventType `storm:"index"`

	StartTime string
	EndTime   string
}

type Team struct {
}