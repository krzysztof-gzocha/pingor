// +build unit

package metric

import (
	"context"
	"testing"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestInstrumentedSuccessRateChecker_Check(t *testing.T) {
	checkerMock := &mock.CheckerMock{}
	ctx := context.TODO()
	innerResult := result.DefaultResult{}
	checkerMock.On("Check", ctx).Once().Return(innerResult)

	checker := NewInstrumentedSuccessRateChecker(checkerMock)
	mainResult := checker.Check(ctx)

	assert.Equal(t, mainResult, innerResult)
	checkerMock.AssertExpectations(t)
}
