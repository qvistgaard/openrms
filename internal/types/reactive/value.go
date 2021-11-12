package reactive

import (
	"context"
	"errors"
	"github.com/reactivex/rxgo/v2"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sync"
	"time"
)

type ValuePostProcessor func(observable rxgo.Observable)
type ValueModifierFunc func(interface{}) interface{}

type Owner struct {
	Type string
	Id   interface{}
}

type Value struct {
	channel     chan rxgo.Item
	observable  rxgo.Observable
	valueType   reflect.Type
	baseValue   interface{}
	value       interface{}
	Annotations Annotations
	lock        *sync.Mutex
	locked      bool
	modifiers   []ValueModifierFunc
}

func (r *Value) RegisterObserver(observer func(rxgo.Observable)) {
	observer(r.observable)
}

func (r *Value) Modifier(modifier ValueModifierFunc) {
	// modifier.register(r)
	r.modifiers = append(r.modifiers, modifier)
}

func (r *Value) Update() error {
	r.value = r.baseValue
	for _, modifier := range r.modifiers {
		r.value = modifier(r.value)
	}
	r.channel <- rxgo.Of(r.value)
	return nil
}

func (r *Value) Set(i interface{}) error {
	of := reflect.TypeOf(i)
	if of != r.valueType {
		log.WithField("type", of.Name()).
			WithField("initial-type", r.valueType.Name()).
			Errorf("invalid baseValue")
		return errors.New("invalid baseValue type")
	}
	r.baseValue = i
	r.Update()
	return nil
}
func (r *Value) Init(ctx context.Context, postProcess ValuePostProcessor) (context.Context, rxgo.Disposable) {
	postProcess(r.observable.Map(func(ctx context.Context, i interface{}) (interface{}, error) {
		return ValueChange{
			Value:       i,
			Type:        r.valueType,
			Annotations: r.Annotations,
			Timestamp:   time.Now(),
		}, nil
	}))
	return r.observable.Connect(ctx)
}

type Annotations map[string]interface{}
type ValueChange struct {
	Value       interface{}
	Type        reflect.Type
	Annotations Annotations
	Timestamp   time.Time
}

func NewDistinctValue(initial interface{}, annotations ...Annotations) Value {
	distinctValueFunc := func(ctx context.Context, i interface{}) (interface{}, error) {
		return i, nil
	}
	value := NewValue(initial, annotations...)
	value.observable = value.observable.DistinctUntilChanged(distinctValueFunc)
	return value
}

func NewValue(initial interface{}, annotations ...Annotations) Value {
	channel := make(chan rxgo.Item)

	mergedAnnotations := Annotations{}
	for _, i := range annotations {
		for k, v := range i {
			mergedAnnotations[k] = v
		}
	}

	value := Value{
		Annotations: mergedAnnotations,
		channel:     channel,
		baseValue:   initial,
		locked:      false,
		lock:        new(sync.Mutex),
		valueType:   reflect.TypeOf(initial),
		observable:  rxgo.FromChannel(channel, rxgo.WithPublishStrategy()),
	}
	return value
}
