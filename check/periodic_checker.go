package check

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/event"
)

// ConnectionCheckEventName will be used as event name when dispatching information about new connection check result
const ConnectionCheckEventName = "connection.check"

// PeriodicCheckerWrapper will use provided internal checker periodically and will trigger an event to dispatcher about each result.
// Period will be getting longer (up to maximalCheckingPeriod) if connection is stable. If any error will be detected, the period
// will go back to minimalCheckingPeriod. It will check success rate and result's time against provided threshold to check if connection is ok or not.
type PeriodicCheckerWrapper struct {
	eventDispatcher       event.DispatcherInterface
	checker               CheckerInterface
	successRateThreshold  float32
	successTimeThreshold  time.Duration
	minimalCheckingPeriod time.Duration
	maximalCheckingPeriod time.Duration
}

// NewPeriodicCheckerWrapper will return new PeriodicCheckerWrapper
func NewPeriodicCheckerWrapper(
	eventDispatcher event.DispatcherInterface,
	checker CheckerInterface,
	successRateThreshold float32,
	successTimeThreshold,
	minimalCheckingPeriod,
	maximalCheckingPeriod time.Duration,
) PeriodicCheckerWrapper {
	return PeriodicCheckerWrapper{
		eventDispatcher:       eventDispatcher,
		checker:               checker,
		successRateThreshold:  successRateThreshold,
		successTimeThreshold:  successTimeThreshold,
		minimalCheckingPeriod: minimalCheckingPeriod,
		maximalCheckingPeriod: maximalCheckingPeriod,
	}
}

// Check should be used to actually start checking process. In order to kill it, you have to kill it's context.
func (c PeriodicCheckerWrapper) Check(ctx context.Context) ResultInterface {
	currentPeriod := c.minimalCheckingPeriod
	for {
		logrus.Debugf("Waiting for %s before next check", currentPeriod)
		select {
		case <-ctx.Done():
			logrus.Debug("PeriodicChecker: exit")
			return Result{}
		case <-time.After(currentPeriod):
		}

		result := c.checker.Check(ctx)
		c.eventDispatcher.Dispatch(ConnectionCheckEventName, result)
		currentPeriod = c.newPeriod(currentPeriod, result)
	}

	return Result{}
}

func (c PeriodicCheckerWrapper) newPeriod(currentPeriod time.Duration, result ResultInterface) time.Duration {
	if result.GetSuccessRate() < c.successRateThreshold {
		return c.minimalCheckingPeriod
	}

	if result.GetTime() > c.successTimeThreshold {
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
