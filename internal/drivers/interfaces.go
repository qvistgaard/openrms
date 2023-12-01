// Package drivers provides the necessary interfaces for implementing device drivers
// in the OpenRMS (Open Racing Management System). These interfaces define the required
// functionalities for cars, tracks, races, and drivers within OpenRMS, ensuring
// standardized interactions and controls across different hardware components.
package drivers

import "github.com/qvistgaard/openrms/internal/types"

// Car interface must be implemented by any device driver representing a car in OpenRMS.
// It defines methods for setting various speed parameters and retrieving the car's identifier.
type Car interface {
	// SetMaxSpeed configures the car's maximum speed, a key parameter for race regulations
	// and safety in OpenRMS.
	SetMaxSpeed(percent uint8)

	// SetPitLaneMaxSpeed sets the maximum speed of the car within the pit lane, crucial for
	// adhering to pit lane regulations within OpenRMS.
	SetPitLaneMaxSpeed(percent uint8)

	// SetMaxBreaking defines the car's braking capacity, essential for driver safety
	// and vehicle control within OpenRMS races.
	SetMaxBreaking(percent uint8)

	// SetMinSpeed establishes a minimum speed for the car, necessary for maintaining
	// consistent race conditions in OpenRMS.
	SetMinSpeed(percent uint8)

	// Id returns the unique identifier for the car, crucial for tracking and management
	// in the OpenRMS.
	Id() types.Id
}

// PitLaneLapCounting enum provides modes for lap counting within the pit lane in OpenRMS.
// This allows flexible race management and strategy development.
type PitLaneLapCounting int

const (
	// LapCountingOnEntry indicates lap counting as a car enters the pit lane.
	LapCountingOnEntry PitLaneLapCounting = iota

	// LapCountingOnExit indicates lap counting on a car's exit from the pit lane.
	LapCountingOnExit
)

// PitLane interface needs to be implemented for managing the pit lane in OpenRMS.
// It allows configuration of lap counting.
type PitLane interface {
	// LapCounting sets up lap counting in the pit lane.
	LapCounting(enabled bool, option PitLaneLapCounting)
}

// Track interface should be implemented by device drivers representing a race track in OpenRMS.
// It provides methods to set track-wide parameters and manage pit lane features.
type Track interface {
	// MaxSpeed sets the track's speed limit, a critical safety and regulation aspect
	// in OpenRMS.
	MaxSpeed(percent uint8)

	// PitLane grants access to the PitLane interface, providing detailed control over
	// pit lane operations within OpenRMS.
	PitLane() PitLane
}

// Race interface is required for device drivers that manage racing events in OpenRMS.
// It includes methods for controlling the various stages and states of a race.
type Race interface {
	// Start triggers the beginning of a race in OpenRMS, marking the commencement
	// of the competitive event.
	Start()

	// Flag signals specific conditions or events during the race.
	Flag()

	// Pause temporarily halts the race, an important feature for managing unforeseen
	// situations in OpenRMS.
	Pause()

	// Stop ends the race.
	Stop()
}

// Driver interface is crucial for any device driver handling driver interactions in OpenRMS.
// It defines methods for starting and stopping driver interactions, and for accessing
// cars, tracks, and races.
type Driver interface {
	// Start initiates the driver's interaction with the race, facilitating event
	// monitoring and control in OpenRMS.
	Start(chan<- Event) error

	// Stop ends the driver's participation in the race, ceasing all active interactions
	// in OpenRMS.
	Stop() error

	// Car retrieves the Car interface for a specified car by its Id, key for
	// driver-car interaction in OpenRMS.
	Car(car types.Id) Car

	// Track provides access to the Track interface, essential for navigating
	// and interacting with the race environment in OpenRMS.
	Track() Track

	// Race offers the Race interface, enabling control and management of race-specific
	// actions and states within OpenRMS.
	Race() Race
}
