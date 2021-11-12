package postprocess

import (
	ctx "context"
	"github.com/qvistgaard/openrms/internal/types/reactive"
	"github.com/reactivex/rxgo/v2"
	"sync"
)

type PostProcessor interface {
	Configure(observable rxgo.Observable)
}

type PostProcess struct {
	postProcessors []PostProcessor
	waitGroup      sync.WaitGroup
	channel        chan rxgo.Item
	observable     rxgo.Observable
}

func CreatePostProcess(postProcessors []PostProcessor) *PostProcess {
	channel := make(chan rxgo.Item)

	process := &PostProcess{
		postProcessors: postProcessors,
		channel:        channel,
		waitGroup:      sync.WaitGroup{},
		observable:     rxgo.FromChannel(channel, rxgo.WithPublishStrategy()),
	}

	for _, pp := range postProcessors {
		pp.Configure(process.observable)
	}
	return process
}

func (p *PostProcess) ValuePostProcessor() reactive.ValuePostProcessor {
	return func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			p.channel <- rxgo.Of(i)
		})
	}
}

func (p *PostProcess) Init(context ctx.Context) {
	p.observable.Connect(context)
}

/*
import (
	"github.com/qvistgaard/openrms/internal/state"
	"sync"
)

type PostProcessor interface {
	CarChannel() chan<- state.CarState
	RaceChannel() chan<- state.CourseState
	RunServer()
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
*/
