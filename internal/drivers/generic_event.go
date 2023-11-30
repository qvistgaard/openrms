package drivers

func GenericEvent(car Car) Event {
	return genericEvent{car}
}

type genericEvent struct {
	car Car
}

func (g genericEvent) Car() Car {
	return g.car
}
