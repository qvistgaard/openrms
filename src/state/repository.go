package state

type Repository interface {
	Get(n string) StateInterface
	All() map[string]StateInterface
	Reset()
	ResetChanges()
	Changes() map[string]StateInterface
}

func CreateInMemoryRepository() *InMemoryRepository {
	i := new(InMemoryRepository)
	i.state = make(map[string]StateInterface)
	return i
}

type InMemoryRepository struct {
	state map[string]StateInterface
}

func (r *InMemoryRepository) All() map[string]StateInterface {
	return r.state
}

func (r *InMemoryRepository) Get(n string) StateInterface {
	if val, ok := r.state[n]; ok {
		return val
	} else {
		r.state[n] = CreateState(r, n, nil)
		r.state[n].initialize()
	}
	return r.state[n]
}

func (r *InMemoryRepository) Reset() {
	state := make(map[string]StateInterface)
	for n, element := range r.state {
		r.state[n] = CreateState(r, n, element.Initial())
	}
	r.state = state
}

func (r *InMemoryRepository) ResetChanges() {
	for _, element := range r.state {
		element.reset()
	}
}

func (r *InMemoryRepository) Changes() map[string]StateInterface {
	var changes = make(map[string]StateInterface)
	for key, element := range r.state {
		if element.Changed() {
			changes[key] = element
		}
	}
	return changes
}
