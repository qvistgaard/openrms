package main

import "github.com/google/uuid"

type Race struct {
	uuid    uuid.UUID
	drivers []Driver
}

type RaceInterface interface {
}

type Practice struct {
	finished bool
	started  bool
}

func main() {

}

type Driver struct {
	id         byte
	car        Car
	controller Controller
}

type Controller struct {
	paired bool
}

type Car struct {
}

type State interface {
	name()
}
