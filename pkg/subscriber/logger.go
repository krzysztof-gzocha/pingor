package subscriber

import (
	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
)

// LogConnectionCheckResult is subscriber that can be used to print basic connection check results in the console
func LogConnectionCheckResult(arg interface{}) {
	res, ok := arg.(result.ResultInterface)
	if !ok {
		return
	}

	logrus.
		WithField("time", res.GetTime().String()).
		WithField("successRate", res.GetSuccessRate()*100).
		Infof("Connection check")
}
