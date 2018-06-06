package multiple

import (
	"context"
	"time"

	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/pkg/check"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
)

// Checker is able to aggregate multiple checkers and treat them as single checker.
// Result will be combination of all the results from individual checkers
type Checker struct {
	checkers             []check.CheckerInterface
	successRateThreshold float32
	successTimeThreshold time.Duration
	singleCheckTimeout   time.Duration
}

// NewChecker will return instance of Checker
func NewChecker(
	singleCheckTimeout time.Duration,
	successRateThreshold float32,
	successTimeThreshold time.Duration,
	checkers ...check.CheckerInterface,
) Checker {
	return Checker{
		singleCheckTimeout:   singleCheckTimeout,
		successRateThreshold: successRateThreshold,
		successTimeThreshold: successTimeThreshold,
		checkers:             checkers,
	}
}

// Check will run all the checkers, combine their's results into single result and return it.
// Each checker will run in separate go-routine
func (c Checker) Check(ctx context.Context) result.ResultInterface {
	overallResult := result.Result{Success: true}
	var wg sync.WaitGroup
	for _, checker := range c.checkers {
		wg.Add(1)
		go func(ctx context.Context, checker check.CheckerInterface, wg *sync.WaitGroup, overallResult *result.Result) {
			logrus.Debugf("Starting checker: %T", checker)
			wrappedCtx, cancelFunc := context.WithTimeout(ctx, c.singleCheckTimeout)
			singleResult := checker.Check(wrappedCtx)
			cancelFunc()
			if !singleResult.IsSuccess() {
				overallResult.Success = false
			}
			overallResult.SubResults = append(overallResult.SubResults, singleResult)
			logrus.Debugf("Checker %T is done", checker)
			wg.Done()
		}(ctx, checker, &wg, &overallResult)
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
