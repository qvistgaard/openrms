package state

type Rule interface {
	// Handle(race connector.Connector, telemetry queue.Queue, car *Car, event ipc.Event)
	InitializeCarState(car *Car)
}
