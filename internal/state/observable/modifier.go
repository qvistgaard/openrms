package observable

// modifier is an internal structure to represent a modifier function with priority.
type modifier[T any] struct {
	modifier Modifier[T]
	priority int
}

// Modifier represents a function type that can modify an observable value and indicate whether the modification is enabled.
type Modifier[T any] func(T) (T, bool)
