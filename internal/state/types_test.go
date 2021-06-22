package state

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"math"
	"os"
	"testing"
	"time"
)

func TestPercentageToUint8(t *testing.T) {
	math.Round(float64(Percent(100)))
	assert.Equal(t, uint8(255), PercentToUint8(float64(Percent(100))))
	assert.Equal(t, uint8(127), PercentToUint8(float64(Percent(50))))
	assert.Equal(t, uint8(25), PercentToUint8(float64(Percent(10))))
}

func TestUint8ToPercentage(t *testing.T) {
	assert.Equal(t, Percent(100), PercentFromUint8(255))
	assert.Equal(t, Percent(50), PercentFromUint8(127))
	assert.Equal(t, Percent(10), PercentFromUint8(25))
}

func Test(t *testing.T) {
	f, err := os.Open("audio/pit-now.mp3")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	// n := format.SampleRate.N(time.Second / 18)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}
