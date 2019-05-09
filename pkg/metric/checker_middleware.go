package metric

import (
	"context"

	"github.com/krzysztof-gzocha/pingor/pkg/check"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
)

type InstrumentedChecker struct {
	checker check.CheckerInterface
}

func NewInstrumentedChecker(checker check.CheckerInterface) *InstrumentedChecker {
	return &InstrumentedChecker{
		checker: checker,
	}
}

func (i *InstrumentedChecker) Check(ctx context.Context) result.ResultInterface {
	checkResult := i.checker.Check(ctx)
	RegisterResult(checkResult)

	return checkResult
}
