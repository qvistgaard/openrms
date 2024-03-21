package v3

import (
	"github.com/go-echarts/statsview"
	"github.com/go-echarts/statsview/viewer"
	"github.com/qvistgaard/openrms/internal/drivers"
	"github.com/qvistgaard/openrms/internal/drivers/devices/oxigen/serial"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	_ "net/http/pprof"
	"os"
	"runtime/debug"
	"testing"
	"time"
)

func TestDriver3xCommunications(t *testing.T) {
	viewer.SetConfiguration(viewer.WithTheme(viewer.ThemeWesteros), viewer.WithMaxPoints(1000))
	mgr := statsview.New()
	mgr.Register(viewer.NewGoroutinesViewer())
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// Start() runs a HTTP server at `localhost:18066` by default.
	go mgr.Start()

	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	debug.SetMemoryLimit(1907088)

	connection, err := serial.CreateUSBConnection(nil)
	assert.Nil(t, err)
	assert.NotNil(t, connection)

	implement, err := CreateDriver(connection)
	assert.Nil(t, err)
	assert.NotNil(t, implement)

	received := make(chan drivers.Event)

	implement.Start(received)

	// var mem runtime.MemStats
	go func() {
		for {
			select {
			case <-received:

				/*				runtime.ReadMemStats(&mem)
								log.Info().
									Int("channel", len(received)).
									Uint64("objects", mem.HeapObjects).
									Uint64("alloc", mem.HeapAlloc).
									Uint64("mem", mem.TotalAlloc).
									Msg("Data received")*/
			}
		}
	}()

	<-time.After(600 * time.Second)

}
