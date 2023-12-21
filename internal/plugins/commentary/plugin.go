package commentary

type Plugin struct {
}

type Config {

}

func (p Plugin) Priority() int {
	return 10
}

func (p Plugin) Name() string {
	return "commentary"
}
