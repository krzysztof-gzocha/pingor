// +build unit

package mock

import (
	"context"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/stretchr/testify/mock"
)

// CheckerMock is mocking Checker
type CheckerMock struct {
	mock.Mock
	Result result.Result
}

// Check is mocking Checker.Check
func (m CheckerMock) Check(ctx context.Context) result.Result {
	args := m.Called(ctx)
	return args.Get(0).(result.Result)
}
