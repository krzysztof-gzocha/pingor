// +build unit

package check

import (
	"testing"
	"time"

	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewMultipleChecker(t *testing.T) {
	checkerMock := checkerMock{}
	checkerMock.On("Check").Return(Result{})
	checker := NewMultipleChecker(time.Second, checkerMock, checkerMock)

	assert.Implements(t, (*CheckerInterface)(nil), checker)
	assert.Len(t, checker.checkers, 2)
	assert.Equal(t, checker.singleCheckTimeout, time.Second)
}

func TestMultipleChecker_Check(t *testing.T) {
	ctx := context.TODO()
	unsuccessfulResult := Result{Success: false, SuccessRate: 0.3, Time: time.Second}
	successfulResult := Result{Success: true, SuccessRate: 0.8, Time: time.Millisecond * 500}
	successChecker := new(checkerMock)
	unsuccessfulChecker := checkerMock{result: unsuccessfulResult}
	successChecker.
		On("Check", mock.AnythingOfType("*context.timerCtx")).
		Once().
		Return(successfulResult)
	unsuccessfulChecker.
		On("Check", mock.AnythingOfType("*context.timerCtx")).
		Once().
		Return(unsuccessfulResult)

	checker := NewMultipleChecker(time.Second, successChecker, unsuccessfulChecker)
	result := checker.Check(ctx)
	assert.NotNil(t, result)
	assert.False(t, result.IsSuccess())
	assert.Equal(t, float32(0.55), result.GetSuccessRate())
	assert.Equal(t, time.Millisecond*750, result.GetTime())
	assert.True(t, successChecker.AssertExpectations(t))
}

type checkerMock struct {
	mock.Mock
	result ResultInterface
}

func (m checkerMock) Check(ctx context.Context) ResultInterface {
	args := m.Called(ctx)
	return args.Get(0).(ResultInterface)
}
