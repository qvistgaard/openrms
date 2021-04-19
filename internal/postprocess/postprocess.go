package postprocess

import (
	"github.com/qvistgaard/openrms/internal/state"
	"sync"
)

type PostProcessor interface {
	CarChannel() chan<- state.CarChanges
	RaceChannel() chan<- state.CourseChanges
	Process()
}

type PostProcess struct {
	postProcessors []PostProcessor
	carChannel     chan state.Car
	raceChannel    chan state.Course
	waitGroup      sync.WaitGroup
}

func CreatePostProcess(postProcessors []PostProcessor) PostProcess {
	return PostProcess{
		postProcessors: postProcessors,
		carChannel:     make(chan state.Car, 1024),
		raceChannel:    make(chan state.Course, 1024),
		waitGroup:      sync.WaitGroup{},
	}
}

func (p *PostProcess) PostProcessCar(changes state.CarChanges) {
	for _, pp := range p.postProcessors {
		pp.CarChannel() <- changes
	}
}

func (p *PostProcess) PostProcessRace(changes state.CourseChanges) {
	for _, pp := range p.postProcessors {
		pp.RaceChannel() <- changes
	}
}
