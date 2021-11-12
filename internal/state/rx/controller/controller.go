package controller

import (
	"context"
	"github.com/qvistgaard/openrms/internal/types/annotations"
	"github.com/qvistgaard/openrms/internal/types/reactive"
)

type Controller struct {
	triggerValue    *reactive.Percent
	buttonTrackCall *reactive.Boolean
}

func NewController(a reactive.Annotations) *Controller {
	return &Controller{
		triggerValue:    reactive.NewPercentAll(0, a, reactive.Annotations{annotations.CarValueFieldName: "trigger-value"}),
		buttonTrackCall: reactive.NewBoolean(false, a, reactive.Annotations{annotations.CarValueFieldName: "track-call"}),
	}
}

func (c *Controller) TriggerValue() *reactive.Percent {
	return c.triggerValue
}

func (c *Controller) ButtonTrackCall() *reactive.Boolean {
	return c.buttonTrackCall
}

func (c *Controller) Init(ctx context.Context, postProcess reactive.ValuePostProcessor) {
	c.buttonTrackCall.Init(ctx, postProcess)
	c.triggerValue.Init(ctx, postProcess)
}
