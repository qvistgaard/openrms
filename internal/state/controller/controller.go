package controller

import (
	"github.com/qvistgaard/openrms/internal/state/observable"
)

func NewController() *Controller {
	controller := &Controller{}

	controller.initObservableProperties()

	controller.filters()

	return controller
}

func (c *Controller) initObservableProperties() {
	c.triggerValue = observable.Create(uint8(0))
	c.buttonTrackCall = observable.Create(false)
}

func (c *Controller) filters() {
	c.buttonTrackCall.Filter(observable.DistinctBooleanChange())
}

type Controller struct {
	buttonTrackCall *observable.Value[bool]
	triggerValue    *observable.Value[uint8]
}

func (c *Controller) TriggerValue() *observable.Value[uint8] {
	return c.triggerValue
}

func (c *Controller) ButtonTrackCall() *observable.Value[bool] {
	return c.buttonTrackCall
}
