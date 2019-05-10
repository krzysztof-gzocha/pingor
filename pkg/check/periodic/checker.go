package periodic

import (
	"context"
	"time"

	"github.com/krzysztof-gzocha/pingor/pkg/check"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/event"
	"github.com/krzysztof-gzocha/pingor/pkg/log"
)

// ConnectionCheckEventName will be used as event name when dispatching information about new connection check result
const ConnectionCheckEventName = "connection.check"

// Checker will use provided internal checker periodically and will trigger an event to dispatcher about each result.
// Period will be getting longer (up to maximalCheckingPeriod) if connection is stable. If any error will be detected, the period
// will go back to minimalCheckingPeriod. It will check success rate and result's time against provided threshold to check if connection is ok or not.
type Checker struct {
	logger                log.Logger
	eventDispatcher       event.Dispatcher
	checker               check.Checker
	successRateThreshold  float32
	successTimeThreshold  time.Duration
	minimalCheckingPeriod time.Duration
	maximalCheckingPeriod time.Duration
}

// NewChecker will return new Checker
func NewChecker(
	logger log.Logger,
	eventDispatcher event.Dispatcher,
	checker check.Checker,
	minimalCheckingPeriod,
	maximalCheckingPeriod time.Duration,
) Checker {
	return Checker{
		logger:                logger,
		eventDispatcher:       eventDispatcher,
		checker:               checker,
		minimalCheckingPeriod: minimalCheckingPeriod,
		maximalCheckingPeriod: maximalCheckingPeriod,
	}
}

// Check should be used to actually start checking process. In order to kill it, you have to kill it's context.
func (c Checker) Check(ctx context.Context) result.Result {
	c.periodicCheck(ctx)

	return result.DefaultResult{}
}

// periodicCheck will run periodic checks on provided checker.
// It's implemented only to get rid of dead code in Check method
func (c Checker) periodicCheck(ctx context.Context) {
	currentPeriod := c.minimalCheckingPeriod
	for {
		c.logger.
			WithField("period", currentPeriod.String()).
			Debugf("%T: Waiting for %s before next check", c, currentPeriod.String())
		select {
		case <-ctx.Done():
			return
		case <-time.After(currentPeriod):
		}

		res := result.DefaultMeasuredAtResult{Result: c.checker.Check(ctx), MeasuredAt: time.Now()}
		c.eventDispatcher.Dispatch(ConnectionCheckEventName, res)
		currentPeriod = c.newPeriod(currentPeriod, res)
	}
}

func (c Checker) newPeriod(currentPeriod time.Duration, result result.Result) time.Duration {
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
