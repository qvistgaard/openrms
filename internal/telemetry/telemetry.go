package telemetry

type Metric struct {
	Name  string
	Value interface{}
}

type Metrics interface {
	Metrics() []Metric
}
