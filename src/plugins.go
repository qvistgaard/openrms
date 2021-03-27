package main

import "github.com/google/uuid"

type RacePlugin interface {
	id() uuid.UUID
	drivers() []*Driver
}
