package state

const (
	RaceEvent = "race-event"
)

func CreateCar(id uint8, settings map[string]interface{}, rules []Rule) *Car {
	c := new(Car)
	c.id = id
	c.settings = settings
	c.state = make(map[string]StateInterface)

	c.ResetState()
	c.Get(RaceEvent).Set(nil)

	for _, r := range rules {
		r.InitializeCarState(c)
	}
	for _, s := range c.state {
		s.initialize()
	}
	return c
}

type Car struct {
	id       uint8
	settings map[string]interface{}
	state    map[string]StateInterface
}

func (c *Car) Get(n string) StateInterface {
	if val, ok := c.state[n]; ok {
		return val
	} else {
		c.state[n] = createState(c, n, nil)
		c.state[n].initialize()
	}
	return c.state[n]
}

func (c *Car) State() map[string]StateInterface {
	return c.state
}

func (c *Car) ResetState() {
	c.state = make(map[string]StateInterface)
	for n, element := range c.settings {
		c.state[n] = createState(c, n, element)
	}
}

func (c *Car) ResetChanges() {
	for _, element := range c.state {
		element.reset()
	}
}

func (c *Car) StateChanges() map[string]StateInterface {
	var changes = make(map[string]StateInterface)
	for key, element := range c.state {
		if element.Changed() {
			changes[key] = element
		}
	}
	return changes
}

func (c *Car) Id() uint8 {
	return c.id
}
