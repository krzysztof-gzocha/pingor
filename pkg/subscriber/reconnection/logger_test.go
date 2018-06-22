// +build unit

package reconnection

import (
	"testing"

	"time"

	"math"

	"github.com/krzysztof-gzocha/pingor/pkg/check/formatter/json"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber"
	"github.com/stretchr/testify/assert"
)

func TestLogger_LogReconnection_BadEvent(t *testing.T) {
	assert.NotPanics(t, func() {
		NewLogger(json.Formatter).LogReconnection(struct{}{})
	})
}

func TestLogger_LogReconnection_CorrectEvent(t *testing.T) {
	res := result.TimeResult{
		MeasuredAt: time.Now(),
		Result:     result.Result{SuccessRate: 1, Success: true},
	}
	assert.NotPanics(t, func() {
		NewLogger(json.Formatter).LogReconnection(subscriber.ReconnectionEvent{
			LastSuccess:         res,
			FirstConnectionDrop: res,
			LastConnectionDrop:  res,
			CurrentResult:       res,
		})
	})
}

func TestLogger_LogReconnection_FormatterError(t *testing.T) {
	res := result.TimeResult{
		MeasuredAt: time.Now(),
		Result:     result.Result{SuccessRate: float32(math.NaN()), Success: true},
	}
	assert.NotPanics(t, func() {
		NewLogger(json.Formatter).LogReconnection(subscriber.ReconnectionEvent{
			LastSuccess:         res,
			FirstConnectionDrop: res,
			LastConnectionDrop:  res,
			CurrentResult:       res,
		})
	})
}
