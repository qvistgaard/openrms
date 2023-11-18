package race

/*
import (
	"github.com/qvistgaard/openrms/internal/config/application"
	"github.com/qvistgaard/openrms/internal/state"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestRule_Notify(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := &application.Context{
		Rules: state.CreateRuleList(),
	}
	ctx.Rules.Append(&Rule{
		ready: make(map[state.CarId]*state.Car),
	})
	ctx.Course = state.CreateCourse(&state.CourseConfig{}, ctx.Rules)
	ctx.Course.Set(state.RMSStatus, state.Initialized)

	c := state.CreateCar(1, nil, ctx.Rules)

	ctx.Course.Set(ActiveView, Started)

	c.Set(state.ControllerBtnTrackCall, true)

	time.Sleep(20 * time.Second)
}
*/
