package reactive

import (
	"context"
	"errors"
	"github.com/reactivex/rxgo/v2"
	log "github.com/sirupsen/logrus"
	"reflect"
	"sort"
	"sync"
	"time"
)

type Condition int

const (
	NoCondition Condition = iota
	IfLessThen
	IfGreaterThen
)

type ValuePostProcessor func(observable rxgo.Observable)
type ValueModifierFunc func(interface{}) (interface{}, bool)
type ValueModifier struct {
	Func     ValueModifierFunc
	Priority int
}

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
	modifiers   []*ValueModifier
}

func (r *Value) RegisterObserver(observer func(rxgo.Observable)) {
	observer(r.observable)
}

func (r *Value) Modifier(modifier ValueModifierFunc, priority int) {
	r.modifiers = append(r.modifiers, &ValueModifier{
		Func:     modifier,
		Priority: priority,
	})

	sort.Slice(r.modifiers, func(i, j int) bool {
		if r.modifiers[i].Priority > r.modifiers[j].Priority {
			return true
		} else {
			return false
		}
	})
}

func (r *Value) Update() error {
	r.value = r.baseValue
	for _, modifier := range r.modifiers {
		if v, enabled := modifier.Func(r.value); enabled {
			r.value = v
		}
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

func NewDistinctValueFunc(initial interface{}, distinctFunc func(ctx context.Context, i interface{}) (interface{}, error), annotations ...Annotations) Value {
	value := NewValue(initial, annotations...)
	value.observable = value.observable.DistinctUntilChanged(distinctFunc)
	return value
}

func NewDistinctValue(initial interface{}, annotations ...Annotations) Value {
	distinctValueFunc := func(ctx context.Context, i interface{}) (interface{}, error) {
		return i, nil
	}
	return NewDistinctValueFunc(initial, distinctValueFunc, annotations...)
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
