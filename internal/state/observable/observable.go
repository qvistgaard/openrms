package observable

import (
	"sort"
)

// modifier is an internal structure to represent a modifier function with priority.
type modifier[T any] struct {
	modifier Modifier[T]
	priority int
}

// Modifier represents a function type that can modify an observable value and indicate whether the modification is enabled.
type Modifier[T any] func(T) (T, bool)

// Filter represents a function type used for filtering observable value changes.
type Filter[T any] func(T, T) bool

// Observer represents a function type used for observing changes in an observable value along with associated annotations.
type Observer[T any] func(T, Annotations)

// Observable is a data structure that allows observing and modifying a value.
type Observable[T any] struct {
	baseValue   T
	value       T
	modifiers   []modifier[T]
	observers   []Observer[T]
	annotations Annotations
	filters     []Filter[T]
}

// Create creates a new Observable with the specified initial value and annotations.
// It returns a pointer to the created Observable.
func Create[T any](initialValue T, annotations ...Annotation) *Observable[T] {
	mergedAnnotations := Annotations{}
	for _, v := range annotations {
		mergedAnnotations[v.Key] = v
	}

	return &Observable[T]{
		baseValue:   initialValue,
		value:       initialValue,
		modifiers:   make([]modifier[T], 0),
		observers:   make([]Observer[T], 0),
		filters:     make([]Filter[T], 0),
		annotations: mergedAnnotations,
	}
}

// RegisterObserver registers an observer function to be notified when the value changes.
// It appends the observer to the list of observers.
// Returns a pointer to the Observable for method chaining.
func (o *Observable[T]) RegisterObserver(observer Observer[T]) *Observable[T] {
	o.observers = append(o.observers, observer)
	return o
}

// Filter registers a filter function to control whether changes to the observable value should be published.
// It appends the filter to the list of filters.
// Returns a pointer to the Observable for method chaining.
func (o *Observable[T]) Filter(filter Filter[T]) *Observable[T] {
	o.filters = append(o.filters, filter)
	return o
}

// Modifier adds a modifier function with a specified priority.
// Modifier functions can transform the value, and their order of execution is determined by priority.
// Modifiers with higher priority values execute first.
// It appends the modifier to the list of modifiers and sorts them by priority.
// Returns a pointer to the Observable for method chaining.
func (o *Observable[T]) Modifier(fn Modifier[T], priority int) *Observable[T] {
	o.modifiers = append(o.modifiers, modifier[T]{
		modifier: fn,
		priority: priority,
	})
	sort.Slice(o.modifiers, func(i, j int) bool {
		if o.modifiers[i].priority > o.modifiers[j].priority {
			return true
		} else {
			return false
		}
	})
	return o
}

// Set sets a new value for the Observable. It applies all registered modifiers
// updates the current value, and notifies observers if the change is verified by filters.
func (o *Observable[T]) Set(value T) {
	baseValue := value
	for _, modifier := range o.modifiers {
		if v, enabled := modifier.modifier(value); enabled {
			value = v
		}
	}

	if o.publishChange(o.value, value) {
		o.baseValue = baseValue
		o.value = value
		for _, observer := range o.observers {
			observer(o.value, o.annotations)
		}
	}
}

// publishChange checks whether the change from 'old' to 'new' value should be published based on registered filters.
// It returns true if the change is allowed by all filters; otherwise, it returns false.
func (o *Observable[T]) publishChange(old T, new T) bool {
	for _, filter := range o.filters {
		if !filter(old, new) {
			return false
		}
	}
	return true
}
