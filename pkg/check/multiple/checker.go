package multiple

import (
	"context"
	"time"

	"sync"

	"github.com/krzysztof-gzocha/pingor/pkg/check"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/log"
)

// Checker is able to aggregate multiple checkers and treat them as single checker.
// DefaultResult will be combination of all the results from individual checkers
type Checker struct {
	logger               log.Logger
	checkers             []check.Checker
	successRateThreshold float32
	successTimeThreshold time.Duration
	singleCheckTimeout   time.Duration
}

// NewChecker will return instance of Checker
func NewChecker(
	logger log.Logger,
	singleCheckTimeout time.Duration,
	successRateThreshold float32,
	successTimeThreshold time.Duration,
	checkers ...check.Checker,
) Checker {
	return Checker{
		logger:               logger,
		singleCheckTimeout:   singleCheckTimeout,
		successRateThreshold: successRateThreshold,
		successTimeThreshold: successTimeThreshold,
		checkers:             checkers,
	}
}

// Check will run all the checkers, combine their's results into single result and return it.
// Each checker will run in separate go-routine
func (c Checker) Check(ctx context.Context) result.Result {
	overallResult := result.DefaultResult{Success: true}
	var wg sync.WaitGroup
	for _, checker := range c.checkers {
		wg.Add(1)
		go c.singleCheck(ctx, checker, &wg, &overallResult)
	}
	wg.Wait()

	var totalTime time.Duration
	if len(overallResult.SubResults) > 0 {
		var successRate float32
		for _, subRequest := range overallResult.SubResults {
			successRate += subRequest.GetSuccessRate()
			totalTime += subRequest.GetTime()
		}
		overallResult.SuccessRate = float32(successRate / float32(len(overallResult.SubResults)))
		overallResult.Time = totalTime / time.Duration(len(overallResult.SubResults))
	}

	if overallResult.GetSuccessRate() < c.successRateThreshold || overallResult.GetTime() > c.successTimeThreshold {
		overallResult.Success = false
	}

	return overallResult
}

func (c Checker) singleCheck(
	ctx context.Context,
	checker check.Checker,
	wg *sync.WaitGroup,
	overallResult *result.DefaultResult,
) {
	c.logger.Debugf("%T: Starting checker: %T", c, checker)
	wrappedCtx, cancelFunc := context.WithTimeout(ctx, c.singleCheckTimeout)
	singleResult := checker.Check(wrappedCtx)
	cancelFunc()
	if !singleResult.IsSuccess() {
		overallResult.Success = false
	}
	overallResult.SubResults = append(overallResult.SubResults, singleResult)
	c.logger.Debugf("%T: Checker %T is done", c, checker)
	wg.Done()
}
