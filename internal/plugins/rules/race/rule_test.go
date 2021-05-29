package race

import (
	"github.com/qvistgaard/openrms/internal/config/context"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestRule_Notify(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := &context.Context{
		Rules: state.CreateRuleList(),
	}
	ctx.Rules.Append(&Rule{
		ready: make(map[state.CarId]bool),
	})
	ctx.Course = state.CreateCourse(&state.CourseConfig{}, ctx.Rules)
	ctx.Course.Set(state.RMSStatus, state.Initialized)

	c := state.CreateCar(1, nil, nil, ctx.Rules)

	ctx.Course.Set(State, Started)

	c.Set(state.ControllerBtnTrackCall, true)

	time.Sleep(20 * time.Second)
}
