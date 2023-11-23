package controller

import (
	"context"
	"github.com/qvistgaard/openrms/internal/state/observable"
	annotations2 "github.com/qvistgaard/openrms/internal/types/annotations"
)

type Controller struct {
	buttonTrackCall *observable.Value[bool]
	triggerValue    *observable.Value[uint8]
}

func NewController(annotations ...observable.Annotation) *Controller {
	return &Controller{
		triggerValue:    observable.Create(uint8(0), append(annotations, observable.Annotation{annotations2.CarValueFieldName, "trigger-value"})...),
		buttonTrackCall: observable.Create(false, annotations...),
	}
}

func (c *Controller) TriggerValue() *observable.Value[uint8] {
	return c.triggerValue
}

func (c *Controller) ButtonTrackCall() *observable.Value[bool] {
	return c.buttonTrackCall
}

func (c *Controller) Init(ctx context.Context) {
	// c.buttonTrackCall.Init(ctx)
	// c.triggerValue.Init(ctx)
}
