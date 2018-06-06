package periodic

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/pkg/check"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/event"
)

// ConnectionCheckEventName will be used as event name when dispatching information about new connection check result
const ConnectionCheckEventName = "connection.check"

// Checker will use provided internal checker periodically and will trigger an event to dispatcher about each result.
// Period will be getting longer (up to maximalCheckingPeriod) if connection is stable. If any error will be detected, the period
// will go back to minimalCheckingPeriod. It will check success rate and result's time against provided threshold to check if connection is ok or not.
type Checker struct {
	eventDispatcher       event.DispatcherInterface
	checker               check.CheckerInterface
	successRateThreshold  float32
	successTimeThreshold  time.Duration
	minimalCheckingPeriod time.Duration
	maximalCheckingPeriod time.Duration
}

// NewChecker will return new Checker
func NewChecker(
	eventDispatcher event.DispatcherInterface,
	checker check.CheckerInterface,
	minimalCheckingPeriod,
	maximalCheckingPeriod time.Duration,
) Checker {
	return Checker{
		eventDispatcher:       eventDispatcher,
		checker:               checker,
		minimalCheckingPeriod: minimalCheckingPeriod,
		maximalCheckingPeriod: maximalCheckingPeriod,
	}
}

// Check should be used to actually start checking process. In order to kill it, you have to kill it's context.
func (c Checker) Check(ctx context.Context) result.ResultInterface {
	c.periodicCheck(ctx)

	return result.Result{}
}

// periodicCheck will run periodic checks on provided checker.
// It's implemented only to get rid of dead code in Check method
func (c Checker) periodicCheck(ctx context.Context) {
	currentPeriod := c.minimalCheckingPeriod
	for {
		logrus.Debugf("Waiting for %s before next check", currentPeriod)
		select {
		case <-ctx.Done():
			logrus.Debug("PeriodicChecker: exit")
			return
		case <-time.After(currentPeriod):
		}

		result := c.checker.Check(ctx)
		c.eventDispatcher.Dispatch(ConnectionCheckEventName, result)
		currentPeriod = c.newPeriod(currentPeriod, result)
	}
}

func (c Checker) newPeriod(currentPeriod time.Duration, result result.ResultInterface) time.Duration {
	if !result.IsSuccess() {
		return c.minimalCheckingPeriod
	}

	if currentPeriod >= c.maximalCheckingPeriod {
		currentPeriod = c.maximalCheckingPeriod

		return currentPeriod
	}

	currentPeriod = currentPeriod * 2
	if currentPeriod > c.maximalCheckingPeriod {
		currentPeriod = c.maximalCheckingPeriod
	}

	return currentPeriod
}
