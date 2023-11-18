package car

import (
	"context"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/types/reactive"
)

type Rule interface {

	// ConfigureCarState is responsible for setting up the rule-specific
	// configuration for a Car object. This method is invoked when a new car
	// is detected, allowing each rule to apply its unique configuration to the
	// car's state. This method should not alter the car itself but should
	// prepare any rule-specific settings or prerequisites that are necessary
	// before the rule is initialized.
	//
	// Parameters:
	//   car - A pointer to the Car object representing the state of the car.
	//
	// Note: This method should focus on preparing the rule's configuration
	// and should not depend on the initialization of the rule, as it precedes
	// the InitializeCarState method.
	ConfigureCarState(*car.Car, *reactive.Factory)

	// InitializeCarState is called after all the ConfigureCarState methods
	// from various rules have been executed. This method is responsible for
	// initializing the rule with the prepared configurations. It ensures that
	// the rule is fully set up and ready to interact with the car's state
	// and other rules.
	//
	// Parameters:
	//   car - A pointer to the Car object representing the state of the car.
	//   ctx - Context for the initialization, allowing for operations like
	//         cancellation and passing request-scoped values.
	//   valuePostProcessor - A function or interface used to process and
	//                        initialize new state objects within the rule.
	//                        It allows for the refinement and customization of
	//                        these state objects as per the rule's requirements.
	//
	// Note: This method completes the rule's setup, ensuring that it can
	// effectively interact with the car's state and other rules.
	InitializeCarState(*car.Car, context.Context)
}
