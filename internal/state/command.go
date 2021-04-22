package state

type CarCommand struct {
	CarId CarId
	Name  string
	Value interface{}
}

type CourseCommand struct {
	Name  string
	Value interface{}
}
