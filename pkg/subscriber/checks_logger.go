package subscriber

import (
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/log"
)

// ChecksLogger holds methods to log information about connection checks
type ChecksLogger struct {
	logger log.LoggerInterface
}

// NewChecksLogger will return ChecksLogger
func NewChecksLogger(logger log.LoggerInterface) *ChecksLogger {
	return &ChecksLogger{logger: logger}
}

// LogConnectionCheckResult is subscriber that can be used to print basic connection check results in the console
func (c *ChecksLogger) LogConnectionCheckResult(arg interface{}) {
	res, ok := arg.(result.ResultInterface)
	if !ok {
		return
	}

	c.logger.
		WithField("time", res.GetTime().String()).
		WithField("successRate", res.GetSuccessRate()*100).
		Infof("Connection check")
}
