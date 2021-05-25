package postprocess

import (
	"github.com/qvistgaard/openrms/internal/state"
	"sync"
)

type PostProcessor interface {
	CarChannel() chan<- state.CarState
	RaceChannel() chan<- state.CourseState
	Process()
}

type CommandEmitter interface {
	CommandChannel(chan<- interface{})
}

type PostProcess struct {
	postProcessors []PostProcessor
	carChannel     chan state.Car
	raceChannel    chan state.Course
	waitGroup      sync.WaitGroup
	CommandChannel chan interface{}
}

func CreatePostProcess(postProcessors []PostProcessor) PostProcess {
	commandChannel := make(chan interface{}, 1024)
	for _, v := range postProcessors {
		if ce, ok := v.(CommandEmitter); ok {
			ce.CommandChannel(commandChannel)
		}
	}
	return PostProcess{
		postProcessors: postProcessors,
		carChannel:     make(chan state.Car, 1024),
		raceChannel:    make(chan state.Course, 1024),
		CommandChannel: commandChannel,
		waitGroup:      sync.WaitGroup{},
	}
}

func (p *PostProcess) PostProcessCar(changes state.CarState) {
	for _, pp := range p.postProcessors {
		pp.CarChannel() <- changes
	}
}

func (p *PostProcess) PostProcessRace(changes state.CourseState) {
	for _, pp := range p.postProcessors {
		pp.RaceChannel() <- changes
	}
}
