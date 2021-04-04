package state

const (
	CarEvent        = "car-event"
	CarPitLaneSpeed = "car-pit-lane-speed"
	CarMinSpeed     = "car-min-speed"
	CarMaxSpeed     = "car-max-speed"
	CarMaxBreaking  = "car-max-breaking"
)

func CreateCar(race *Race, id uint8, settings map[string]interface{}, rules []Rule) *Car {
	c := new(Car)
	c.id = id
	c.race = race
	c.settings = settings
	c.state = CreateInMemoryRepository()

	c.state.Reset()
	// c.state.Get(CarEvent).Set(nil)

	for _, r := range rules {
		r.InitializeCarState(c)
	}
	for _, s := range c.state.All() {
		s.initialize()
	}
	return c
}

type Car struct {
	id       uint8
	settings map[string]interface{}
	state    Repository
	race     *Race
}

func (c *Car) Race() *Race {
	return c.race
}

func (c *Car) State() Repository {
	return c.state
}

func (c *Car) Id() uint8 {
	return c.id
}
