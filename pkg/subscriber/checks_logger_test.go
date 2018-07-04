// +build unit

package subscriber

import (
	"testing"

	"time"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	pkgMock "github.com/krzysztof-gzocha/pingor/pkg/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogConnectionCheckResult_BadResult(t *testing.T) {
	checkLogger := NewChecksLogger(&pkgMock.Logger{})
	assert.NotPanics(t, func() {
		checkLogger.LogConnectionCheckResult(struct{}{})
	})
}

func TestLogConnectionCheckResult_EmptyResult(t *testing.T) {
	logger := &pkgMock.Logger{}
	logger.On("WithField", "time", mock.Anything)
	logger.On("WithField", "successRate", mock.Anything)
	logger.On("Infof", mock.Anything, mock.Anything)

	checkLogger := NewChecksLogger(logger)
	assert.NotPanics(t, func() {
		checkLogger.LogConnectionCheckResult(result.Result{})
	})
	logger.AssertExpectations(t)
}

func TestLogConnectionCheckResult_FullResult(t *testing.T) {
	logger := &pkgMock.Logger{}
	logger.On("WithField", "time", mock.Anything)
	logger.On("WithField", "successRate", mock.Anything)
	logger.On("Infof", mock.Anything, mock.Anything)
	checkLogger := NewChecksLogger(logger)
	assert.NotPanics(t, func() {
		checkLogger.LogConnectionCheckResult(result.Result{
			Time:        time.Second,
			SuccessRate: 75,
		})
	})
	logger.AssertExpectations(t)
}
