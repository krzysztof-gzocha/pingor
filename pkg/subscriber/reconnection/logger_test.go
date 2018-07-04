// +build unit

package reconnection

import (
	"testing"

	"time"

	"math"

	"github.com/krzysztof-gzocha/pingor/pkg/check/formatter/json"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	pkgMock "github.com/krzysztof-gzocha/pingor/pkg/mock"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogger_LogReconnection_BadEvent(t *testing.T) {
	assert.NotPanics(t, func() {
		NewLogger(&pkgMock.Logger{}, json.Formatter).LogReconnection(struct{}{})
	})
}

func TestLogger_LogReconnection_CorrectEvent(t *testing.T) {
	logger := &pkgMock.Logger{}
	logger.On("WithField", "lastSuccessTime", mock.Anything)
	logger.On("WithField", "firstConnectionDrop", mock.Anything)
	logger.On("WithField", "lastConnectionDrop", mock.Anything)
	logger.On("WithField", "current", mock.Anything)
	logger.On("WithField", "lastSuccessRate", mock.Anything)
	logger.On("WithField", "currentSuccessRate", mock.Anything)
	logger.On("WithField", "disconnectionTime", mock.Anything)
	logger.On("Infof", mock.Anything, mock.Anything)
	res := result.TimeResult{
		MeasuredAt: time.Now(),
		Result:     result.Result{SuccessRate: 1, Success: true},
	}
	assert.NotPanics(t, func() {
		NewLogger(logger, json.Formatter).LogReconnection(subscriber.ReconnectionEvent{
			LastSuccess:         res,
			FirstConnectionDrop: res,
			LastConnectionDrop:  res,
			CurrentResult:       res,
		})
	})
}

func TestLogger_LogReconnection_FormatterError(t *testing.T) {
	logger := &pkgMock.Logger{}

	logger.On("WithField", "lastSuccessTime", mock.Anything)
	logger.On("WithField", "firstConnectionDrop", mock.Anything)
	logger.On("WithField", "lastConnectionDrop", mock.Anything)
	logger.On("WithField", "current", mock.Anything)
	logger.On("WithField", "lastSuccessRate", mock.Anything)
	logger.On("WithField", "currentSuccessRate", mock.Anything)
	logger.On("WithField", "disconnectionTime", mock.Anything)
	logger.On("Infof", mock.Anything, mock.Anything)
	res := result.TimeResult{
		MeasuredAt: time.Now(),
		Result:     result.Result{SuccessRate: float32(math.NaN()), Success: true},
	}
	assert.NotPanics(t, func() {
		NewLogger(logger, json.Formatter).LogReconnection(subscriber.ReconnectionEvent{
			LastSuccess:         res,
			FirstConnectionDrop: res,
			LastConnectionDrop:  res,
			CurrentResult:       res,
		})
	})
}
