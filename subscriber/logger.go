package subscriber

import (
	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/check"
)

// LogConnectionCheckResult is subscriber that can be used to print basic connection check results in the console
func LogConnectionCheckResult(arg interface{}) {
	result, ok := arg.(check.ResultInterface)
	if !ok {
		return
	}

	logrus.Infof("Connection check with success rate %.2f%% and time %s", result.GetSuccessRate()*100, result.GetTime())
}
