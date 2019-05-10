// +build unit

package periodic

import (
	"testing"

	"time"

	"context"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/event"
	pkgMock "github.com/krzysztof-gzocha/pingor/pkg/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPeriodicCheckerWrapper(t *testing.T) {
	subChecker := pkgMock.CheckerMock{}
	eventDispatcherMock := eventDispatcherMock{}

	checker := NewChecker(
		&pkgMock.Logger{},
		eventDispatcherMock,
		subChecker,
		time.Second,
		time.Second,
	)

	assert.Equal(t, subChecker, checker.checker)
	assert.Equal(t, eventDispatcherMock, checker.eventDispatcher)
}

func TestPeriodicCheckerWrapper_Check(t *testing.T) {
	logger := &pkgMock.Logger{}
	logger.On("Debugf", mock.Anything, mock.Anything)
	logger.On("WithField", "period", mock.Anything)
	ctx, cancelFunc := context.WithCancel(context.TODO())
	subChecker := pkgMock.CheckerMock{}
	eventDispatcherMock := eventDispatcherMock{}
	subChecker.
		On("Check", ctx).
		Times(3).
		Return(result.DefaultResult{Success: true, SuccessRate: 1})
	eventDispatcherMock.
		On("Dispatch", "connection.check", mock.AnythingOfType("result.DefaultMeasuredAtResult")).
		Times(3)

	checker := NewChecker(
		logger,
		eventDispatcherMock,
		subChecker,
		time.Millisecond*100,
		time.Millisecond*1000,
	)

	go func() {
		<-time.After(time.Millisecond * 800)
		cancelFunc()
	}()

	checker.Check(ctx)
	cancelFunc()

	assert.True(t, subChecker.AssertExpectations(t))
	assert.True(t, eventDispatcherMock.AssertExpectations(t))
	logger.AssertExpectations(t)
}

func TestPeriodicCheckerWrapper_newPeriod(t *testing.T) {
	logger := &pkgMock.Logger{}
	subChecker := pkgMock.CheckerMock{}
	eventDispatcherMock := eventDispatcherMock{}

	minimalPeriod := time.Millisecond * 500
	maximalPeriod := time.Minute
	checker := NewChecker(
		logger,
		eventDispatcherMock,
		subChecker,
		minimalPeriod,
		maximalPeriod,
	)

	scenarios := []struct {
		time     time.Duration
		res      result.DefaultResult
		expected time.Duration
	}{
		{time: time.Second, res: result.DefaultResult{Success: false}, expected: minimalPeriod},
		{time: minimalPeriod * 2, res: result.DefaultResult{Success: false}, expected: minimalPeriod},
		{time: maximalPeriod * 2, res: result.DefaultResult{Success: true}, expected: maximalPeriod},
		{time: maximalPeriod - time.Millisecond, res: result.DefaultResult{Success: true}, expected: maximalPeriod},
	}

	for _, scenario := range scenarios {
		newPeriod := checker.newPeriod(scenario.time, scenario.res)
		assert.Equal(t, scenario.expected, newPeriod)
	}
	logger.AssertExpectations(t)
}

type eventDispatcherMock struct {
	mock.Mock
}

func (m eventDispatcherMock) AttachSubscriber(eventName string, subscriber event.Subscriber) event.Dispatcher {
	m.Called(eventName, subscriber)
	return m
}
func (m eventDispatcherMock) Dispatch(eventName string, arg interface{}) {
	m.Called(eventName, arg)
}
