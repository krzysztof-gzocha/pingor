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
		eventDispatcherMock,
		subChecker,
		time.Second,
		time.Second,
	)

	assert.Equal(t, subChecker, checker.checker)
	assert.Equal(t, eventDispatcherMock, checker.eventDispatcher)
}

func TestPeriodicCheckerWrapper_Check(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	subChecker := pkgMock.CheckerMock{}
	eventDispatcherMock := eventDispatcherMock{}
	subChecker.
		On("Check", ctx).
		Times(3).
		Return(result.Result{Success: true, SuccessRate: 1})
	eventDispatcherMock.
		On("Dispatch", "connection.check", mock.AnythingOfType("result.TimeResult")).
		Times(3)

	checker := NewChecker(
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
}

func TestPeriodicCheckerWrapper_newPeriod(t *testing.T) {
	subChecker := pkgMock.CheckerMock{}
	eventDispatcherMock := eventDispatcherMock{}

	minimalPeriod := time.Millisecond * 500
	maximalPeriod := time.Minute
	checker := NewChecker(
		eventDispatcherMock,
		subChecker,
		minimalPeriod,
		maximalPeriod,
	)

	scenarios := []struct {
		time     time.Duration
		res      result.Result
		expected time.Duration
	}{
		{time: time.Second, res: result.Result{Success: false}, expected: minimalPeriod},
		{time: minimalPeriod * 2, res: result.Result{Success: false}, expected: minimalPeriod},
		{time: maximalPeriod * 2, res: result.Result{Success: true}, expected: maximalPeriod},
		{time: maximalPeriod - time.Millisecond, res: result.Result{Success: true}, expected: maximalPeriod},
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
