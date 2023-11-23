package observable

// Filter represents a function type used for filtering observable value changes.

// Observer represents a function type used for observing changes in an observable value along with associated annotations.
type Observer[T any] func(T, Annotations)

// Observable is an interface that represents an observable value, which can be observed, modified, and filtered.
// It allows for registering observers, applying filters, and modifying the value it holds.
type Observable[T any] interface {

	// RegisterObserver registers an observer function to be notified when the value changes.
	// It takes an Observer[T] and returns the modified Observable[T] for method chaining.
	RegisterObserver(observer Observer[T]) Observable[T]

	// Filter registers a filter function to control whether changes to the observable value should be published.
	// It takes a Filter[T] and returns the modified Observable[T] for method chaining.
	Filter(filter Filter[T]) Observable[T]

	// Modifier adds a modifier function with a specified priority.
	// Modifier functions can transform the value, and their order of execution is determined by priority.
	// Modifiers with higher priority values execute first.
	// It takes a Modifier[T] function and an integer priority, and returns the modified Observable[T] for method chaining.
	Modifier(fn Modifier[T], priority int) Observable[T]

	// Publish notifies registered observers about the latest value change.
	// It triggers the execution of observer functions.
	Publish()

	Update()

	// Set sets a new value for the Observable. It applies all registered modifiers,
	// updates the current value, and notifies observers.
	// It takes a value of type T.
	Set(value T)

	// Get retrieves the current value of the Observable.
	// It returns the current value of type T.
	Get() T
}
