package car

import (
	"context"
	"github.com/qvistgaard/openrms/internal/state/car"
	"github.com/qvistgaard/openrms/internal/types/reactive"
)

type Rule interface {
	ConfigureCarState(*car.Car)
	InitializeCarState(*car.Car, context.Context, reactive.ValuePostProcessor)
}
