package drivers

// Event is an interface representing an event that is associated with a Car.
// It provides a method to retrieve the Car involved in the event.
// Note: The Car returned can be nil, indicating that no specific Car is associated
// with the event or the association is not applicable in the current context.
type Event interface {

	// Car returns the Car associated with the event.
	// It can return nil if there is no Car associated with the event.
	Car() Car
}
