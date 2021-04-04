package state

type Subscriber interface {
	Notify(v *Value)
}
