package reactive

import (
	"context"
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/reactivex/rxgo/v2"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLiterSubtractModifier_Modify(t *testing.T) {

	mod := &LiterSubtractModifier{
		Subtract: 20,
	}
	liter := NewLiter(100)

	output := make(chan ValueChange)

	liter.Init(context.Background(), func(observable rxgo.Observable) {
		observable.DoOnNext(func(i interface{}) {
			output <- i.(ValueChange)
		})
	})

	liter.Modifier(mod, 1)

	liter.Update()
	result := <-output
	assert.Equal(t, types.Percent(80), result.Value)

	mod.Subtract = 30
	liter.Update()
	result = <-output
	assert.Equal(t, types.Percent(70), result.Value)

	log.Error(result)

}
