package events

import "github.com/qvistgaard/openrms/internal/drivers"

type GenericControllerTrackCallButton struct {
	Button
}

func NewControllerTrackCallButton(car drivers.Car, pressed bool) Button {
	return &GenericControllerTrackCallButton{NewGenericButton(car, pressed)}
}

// ControllerTrackCallButton interface represents a specific type of button used
// for controller track call purposes in racing game setups. Currently, it does not
// extend the Button interface with additional methods but serves as a distinct
// type for clarity and potential future expansion.
type ControllerTrackCallButton interface {
	Button
}
