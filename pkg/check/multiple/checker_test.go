// +build unit

package multiple

import (
	"testing"
	"time"

	"context"

	"github.com/krzysztof-gzocha/pingor/pkg/check"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	pkgMock "github.com/krzysztof-gzocha/pingor/pkg/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewMultipleChecker(t *testing.T) {
	checkerMock := pkgMock.CheckerMock{}
	checkerMock.On("Check").Return(result.Result{})
	checker := NewChecker(time.Second, 1, time.Second, checkerMock, checkerMock)

	assert.Implements(t, (*check.CheckerInterface)(nil), checker)
	assert.Len(t, checker.checkers, 2)
	assert.Equal(t, checker.singleCheckTimeout, time.Second)
}

func TestMultipleChecker_Check(t *testing.T) {
	ctx := context.TODO()
	unsuccessfulResult := result.Result{Success: false, SuccessRate: 0.4, Time: time.Millisecond * 40}
	successfulResult := result.Result{Success: true, SuccessRate: 0.6, Time: time.Millisecond * 60}
	successChecker := new(pkgMock.CheckerMock)
	unsuccessfulChecker := pkgMock.CheckerMock{Result: unsuccessfulResult}
	successChecker.
		On("Check", mock.Anything).
		Once().
		Return(successfulResult)
	unsuccessfulChecker.
		On("Check", mock.Anything).
		Once().
		Return(unsuccessfulResult)

	checker := NewChecker(
		time.Second,
		1,
		time.Second*5,
		successChecker,
		unsuccessfulChecker,
	)
	res := checker.Check(ctx)
	assert.NotNil(t, res)
	assert.False(t, res.IsSuccess())
	assert.Equal(t, float32(0.5), res.GetSuccessRate())
	assert.Equal(t, time.Millisecond*50, res.GetTime())
	assert.True(t, successChecker.AssertExpectations(t))
}

func TestMultipleChecker_Check_ThresholdsConditions(t *testing.T) {
	ctx := context.TODO()
	successfulResult := result.Result{
		Success:     true,
		SuccessRate: 0.85,
		Time:        time.Second,
	}
	successChecker := new(pkgMock.CheckerMock)
	successChecker.
		On("Check", mock.AnythingOfType("*context.timerCtx")).
		Once().
		Return(successfulResult)

	checker := NewChecker(
		time.Second,
		1,
		time.Millisecond,
		successChecker,
	)

	res := checker.Check(ctx)
	assert.NotNil(t, res)
	assert.False(t, res.IsSuccess())
	assert.Equal(t, float32(0.85), res.GetSuccessRate())
	assert.Equal(t, time.Second, res.GetTime())
	assert.True(t, successChecker.AssertExpectations(t))
}
