package events

import "github.com/qvistgaard/openrms/internal/drivers"

// Button interface represents a generic button event. It extends the drivers.Event
// interface and adds functionality to determine the button's state.
type Button interface {
	drivers.Event

	// Pressed returns true if the button is currently pressed, false otherwise.
	Pressed() bool
}

// GenericButton struct is an implementation of the Button interface.
// It represents a basic button with a pressed state.
type GenericButton struct {
	drivers.Event
	pressed bool
}

// NewGenericButton is a constructor function for GenericButton.
// It creates and returns a new GenericButton instance.
//
// car: The drivers.Car associated with the button.
// pressed: The pressed state of the button.
//
// returns: An instance of Button, specifically a GenericButton.
func NewGenericButton(car drivers.Car, pressed bool) Button {
	return &GenericButton{NewGenericEvent(car), pressed}
}

// Pressed implements the Pressed method of the Button interface.
// It returns the current pressed state of the button.
func (g GenericButton) Pressed() bool {
	return g.pressed
}
