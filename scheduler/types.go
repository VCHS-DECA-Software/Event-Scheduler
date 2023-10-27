package scheduler

import (
	"Event-Scheduler/components/proto"
)

// decisions
type Output struct {
	Housings []Housing
	Context  ScheduleContext
	Exams    []Exam
}

type Exam struct {
	Start   int
	Student *proto.Student
}

type Housing struct {
	Room   *proto.Room
	Judges []*Judgement
}

type Judgement struct {
	Judge       *proto.Judge
	Assignments []Assignment
}

type Assignment struct {
	Event *proto.Event
	Group []*proto.Student
}
