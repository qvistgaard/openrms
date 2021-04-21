package state

import log "github.com/sirupsen/logrus"

func CreateState(owner Owner, name string, v interface{}) *Value {
	s := new(Value)
	s.value = v
	s.initial = v
	s.changed = false
	s.owner = owner
	s.name = name
	return s
}

type Change struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type Value struct {
	changed     bool
	name        string
	value       interface{}
	previous    interface{}
	initialized bool
	initial     interface{}
	owner       Owner
	subscribers []Subscriber
}

type StateInterface interface {
	Get() interface{}
	GetPrevious() interface{}
	Set(v interface{})
	Name() string
	Owner() Owner
	initialize()
	Initial() interface{}
	Subscribe(s Subscriber)
	Changed() bool
	reset()
}

func (v *Value) Get() interface{} {
	return v.value
}

func (v *Value) GetPrevious() interface{} {
	return v.previous
}

func (v *Value) Set(value interface{}) {
	if v.value != value {
		v.previous = v.value
		v.value = value
		if v.initialized {
			log.Infof("%s=%+v", v.name, value)
			v.changed = true
			for _, s := range v.subscribers {
				s.Notify(v)
			}
		}
	}
}

func (v *Value) initialize() {
	v.initialized = true
}

func (v *Value) Initial() interface{} {
	return v.initial
}

func (v *Value) reset() {
	v.changed = false
}

func (v *Value) Changed() bool {
	return v.changed
}

func (v *Value) Owner() Owner {
	return v.owner
}

func (v *Value) Name() string {
	return v.name
}

func (v *Value) Subscribe(s Subscriber) {
	v.subscribers = append(v.subscribers, s)
}
