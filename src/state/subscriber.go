package state

import "openrms/telemetry"

type Subscriber interface {
	Notify(v *Value)
}

type Telemetry interface {
	TelemetryProcessor(processor telemetry.Processor)
}
