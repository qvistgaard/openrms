package observable

import "sort"

// Value is a data structure that allows observing and modifying a value.
type Value[T any] struct {
	baseValue   T
	value       T
	modifiers   []modifier[T]
	observers   []Observer[T]
	annotations Annotations
	filters     []Filter[T]
}

// Create creates a new Value with the specified initial value and annotations.
// It returns a pointer to the created Value.
func Create[T any](initialValue T, annotations ...Annotation) *Value[T] {
	mergedAnnotations := Annotations{}
	for _, v := range annotations {
		mergedAnnotations[v.Key] = v
	}

	return &Value[T]{
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
// Returns a pointer to the Value for method chaining.
func (o *Value[T]) RegisterObserver(observer Observer[T]) Observable[T] {
	o.observers = append(o.observers, observer)
	return o
}

// Filter registers a filter function to control whether changes to the observable value should be published.
// It appends the filter to the list of filters.
// Returns a pointer to the Value for method chaining.
func (o *Value[T]) Filter(filter Filter[T]) Observable[T] {
	o.filters = append(o.filters, filter)
	return o
}

// Modifier adds a modifier function with a specified priority.
// Modifier functions can transform the value, and their order of execution is determined by priority.
// Modifiers with higher priority values execute first.
// It appends the modifier to the list of modifiers and sorts them by priority.
// Returns a pointer to the Value for method chaining.
func (o *Value[T]) Modifier(fn Modifier[T], priority int) Observable[T] {
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

// Set sets a new value for the Value. It applies all registered modifiers
// updates the current value, and notifies observers if the change is verified by filters.
func (o *Value[T]) Set(value T) {
	baseValue := value
	for _, modifier := range o.modifiers {
		if v, enabled := modifier.modifier(value); enabled {
			value = v
		}
	}

	if o.isValueChanged(value) {
		o.baseValue = baseValue
		o.value = value
		o.Publish()
	}
}

// Get retrieves the current value of the Value.
// It returns the current value of type T.
func (o Value[T]) Get() T {
	return o.value
}

// Publish notifies registered observers about the latest value change in the Value.
// It triggers the execution of observer functions with the current value and associated annotations.
func (o *Value[T]) Publish() {
	for _, observer := range o.observers {
		observer(o.value, o.annotations)
	}
}

// isValueChanged checks whether the change from 'before' to 'after' value should be updated and published based on registered filters.
// It returns true if the change is allowed by all filters; otherwise, it returns false.
func (o *Value[T]) isValueChanged(after T) bool {
	for _, filter := range o.filters {
		if !filter(o.baseValue, after) {
			return false
		}
	}
	return true
}
