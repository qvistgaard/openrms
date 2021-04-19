package car

type Settings struct {
	Id       uint8
	MaxSpeed uint8 `yaml:"max-speed"`
}

type Repository interface {
	GetCarById(uint82 uint8) map[string]interface{}
}
