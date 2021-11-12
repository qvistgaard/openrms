package state

/*type Repository interface {
	Get(n string) StateInterface
	Create(n string, v interface{})
	All() map[string]StateInterface
	SetDefaults()
	ResetChanges()
	Changes() map[string]StateInterface
}

func CreateInMemoryRepository(owner Owner) *InMemoryRepository {
	i := new(InMemoryRepository)
	i.state = make(map[string]StateInterface)
	i.owner = owner
	return i
}

func (r *InMemoryRepository) Create(n string, v interface{}) {
	r.state[n] = CreateState(r.owner, n, v)
	r.state[n].initialize()
}

type InMemoryRepository struct {
	state map[string]StateInterface
	owner Owner
}

func (r *InMemoryRepository) All() map[string]StateInterface {
	return r.state
}

func (r *InMemoryRepository) Get(n string) StateInterface {
	if val, ok := r.state[n]; ok {
		return val
	} else {
		r.Create(n, nil)
	}
	return r.state[n]
}

func (r *InMemoryRepository) SetDefaults() {
	// state := make(map[string]StateInterface)
	for n, element := range r.state {
		r.state[n] = CreateState(r, n, element.Initial())
	}
	// r.state = state
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
*/
