// +build unit

package mock

import (
	"context"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/stretchr/testify/mock"
)

type CheckerMock struct {
	mock.Mock
	Result result.ResultInterface
}

func (m CheckerMock) Check(ctx context.Context) result.ResultInterface {
	args := m.Called(ctx)
	return args.Get(0).(result.ResultInterface)
}
