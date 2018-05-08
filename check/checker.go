package check

import (
	"context"
	"time"

	"sync"

	"github.com/Sirupsen/logrus"
)

// MultipleChecker is able to aggregate multiple checkers and treat them as single checker.
// Result will be combination of all the results from individual checkers
type MultipleChecker struct {
	checkers           []CheckerInterface
	singleCheckTimeout time.Duration
}

// NewMultipleChecker will return instance of MultipleChecker
func NewMultipleChecker(
	singleCheckTimeout time.Duration,
	checkers ...CheckerInterface,
) MultipleChecker {
	return MultipleChecker{
		singleCheckTimeout: singleCheckTimeout,
		checkers:           checkers,
	}
}

// CheckerInterface should be implemented by anything that is able to check current internet connection and return a result
type CheckerInterface interface {
	Check(ctx context.Context) ResultInterface
}

// Check will run all the checkers, combine their's results into single result and return it.
// Each checker will run in separate go-routine
func (c MultipleChecker) Check(ctx context.Context) ResultInterface {
	overallResult := Result{Success: true}
	var wg sync.WaitGroup
	for _, checker := range c.checkers {
		wg.Add(1)
		go func(ctx context.Context, checker CheckerInterface, wg *sync.WaitGroup, overallResult *Result) {
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

	return overallResult
}
