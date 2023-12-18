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
	c.buttonTrackCall = observable.Create(false).Filter(observable.DistinctBooleanChange())
}

func (c *Controller) filters() {
	c.buttonTrackCall.Filter(observable.DistinctBooleanChange())
}

type Controller struct {
	buttonTrackCall observable.Observable[bool]
	triggerValue    observable.Observable[uint8]
}

func (c *Controller) TriggerValue() observable.Observable[uint8] {
	return c.triggerValue
}

func (c *Controller) ButtonTrackCall() observable.Observable[bool] {
	return c.buttonTrackCall
}
