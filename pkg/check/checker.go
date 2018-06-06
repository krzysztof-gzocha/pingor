package check

import (
	"context"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
)

// CheckerInterface should be implemented by anything that is able to check current internet connection and return a result
type CheckerInterface interface {
	Check(ctx context.Context) result.ResultInterface
}
