// +build unit

package check

import (
	"testing"

	"time"

	"context"

	"github.com/krzysztof-gzocha/pingor/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPeriodicCheckerWrapper(t *testing.T) {
	subChecker := checkerMock{}
	eventDispatcherMock := eventDispatcherMock{}

	checker := NewPeriodicCheckerWrapper(
		eventDispatcherMock,
		subChecker,
		0.5,
		time.Second,
		time.Second,
		time.Second,
	)

	assert.Equal(t, subChecker, checker.checker)
	assert.Equal(t, eventDispatcherMock, checker.eventDispatcher)
}

func TestPeriodicCheckerWrapper_Check(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	subChecker := checkerMock{}
	eventDispatcherMock := eventDispatcherMock{}
	subChecker.
		On("Check", ctx).
		Times(3).
		Return(Result{Success: true, SuccessRate: 1})
	eventDispatcherMock.
		On("Dispatch", "connection.check", mock.AnythingOfType("check.Result")).
		Times(3)

	checker := NewPeriodicCheckerWrapper(
		eventDispatcherMock,
		subChecker,
		0.5,
		time.Millisecond*100,
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
}

func TestPeriodicCheckerWrapper_newPeriod(t *testing.T) {
	subChecker := checkerMock{}
	eventDispatcherMock := eventDispatcherMock{}

	minimalPeriod := time.Millisecond * 500
	maximalPeriod := time.Minute
	successTimeThreshold := time.Millisecond * 10
	checker := NewPeriodicCheckerWrapper(
		eventDispatcherMock,
		subChecker,
		0.5,
		successTimeThreshold,
		minimalPeriod,
		maximalPeriod,
	)

	scenarios := []struct {
		time     time.Duration
		res      Result
		expected time.Duration
	}{
		{time: time.Second, res: Result{SuccessRate: 0}, expected: minimalPeriod},
		{time: minimalPeriod * 2, res: Result{SuccessRate: 1, Time: successTimeThreshold * 5}, expected: minimalPeriod},
		{time: maximalPeriod * 2, res: Result{SuccessRate: 1, Time: successTimeThreshold}, expected: maximalPeriod},
		{time: maximalPeriod - time.Millisecond, res: Result{SuccessRate: 1, Time: successTimeThreshold}, expected: maximalPeriod},
	}

	for _, scenario := range scenarios {
		newPeriod := checker.newPeriod(scenario.time, scenario.res)
		assert.Equal(t, scenario.expected, newPeriod)
	}
}

type eventDispatcherMock struct {
	mock.Mock
}

func (m eventDispatcherMock) AttachSubscriber(eventName string, subscriber event.Subscriber) event.DispatcherInterface {
	m.Called(eventName, subscriber)
	return m
}
func (m eventDispatcherMock) Dispatch(eventName string, arg interface{}) {
	m.Called(eventName, arg)
}
