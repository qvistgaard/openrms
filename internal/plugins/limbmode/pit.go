package limbmode

import (
	"github.com/qvistgaard/openrms/internal/state/observable"
	"time"
)

type Sequence struct {
	limbMode observable.Observable[bool]
}

func NewSequence(limbMode observable.Observable[bool]) *Sequence {
	return &Sequence{limbMode: limbMode}
}

func (s *Sequence) Start() error {
	if s.limbMode.Get() == true {
		time.Sleep(10 * time.Second)
	}
	s.limbMode.Set(false)
	return nil

}
