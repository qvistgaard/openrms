package sound

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestAverageFuelLoad(t *testing.T) {
	u := float64(45)
	a := float64(10)
	left := u / a

	log.Info(left)

	if left < 5 {
		log.Info("Announce")
	}

}
