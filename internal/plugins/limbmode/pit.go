package limbmode

import (
	"github.com/qvistgaard/openrms/internal/state/observable"
	log "github.com/sirupsen/logrus"
	"time"
)

type Sequence struct {
	limbMode observable.Observable[bool]
}

func NewSequence(limbMode observable.Observable[bool]) *Sequence {
	return &Sequence{limbMode: limbMode}
}

func (s *Sequence) Start() error {
	log.Info("Limbmode pit sequence started.")

	if s.limbMode.Get() == true {
		log.Info("Limbmode penalty started.")
		time.Sleep(10 * time.Second)
		log.Info("Limbmode penalty completed.")
	}
	s.limbMode.Set(false)
	s.limbMode.Publish()

	log.Info("Limbmode pit sequence completed.")
	return nil

}
