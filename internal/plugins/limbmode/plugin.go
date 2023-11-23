package limbmode

type Plugin struct {
}

func (p Plugin) Priority() int {
	return 10
}
