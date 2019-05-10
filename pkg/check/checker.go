package check

import (
	"context"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
)

// Checker should be implemented by anything that is able to check current internet connection and return a result
type Checker interface {
	Check(ctx context.Context) result.Result
}
