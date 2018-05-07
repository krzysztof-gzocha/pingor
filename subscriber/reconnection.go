package subscriber

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/check"
)

// Reconnection subscriber is responsible to check if connection was re-established. If so it will create proper log about it.
type Reconnection struct {
	previousResult     check.ResultInterface
	lastConnectionDrop time.Time
	printer            check.ResultPrinter
}

// NewReconnectionSubscriber will return a pointer to Reconnection
func NewReconnectionSubscriber(printer check.ResultPrinter) *Reconnection {
	return &Reconnection{printer: printer}
}

// NotifyAboutReconnection is subscriber method that will read the result from provided argument and interpret them.
// In case of when last result had errors and current is clear it will log this information alongside with time details.
func (r *Reconnection) NotifyAboutReconnection(arg interface{}) {
	result, ok := arg.(check.ResultInterface)
	if !ok {
		return
	}

	if r.previousResult == nil {
		r.previousResult = result
		return
	}

	if !r.previousResult.IsSuccess() && result.IsSuccess() {
		logrus.Warnf(
			"Connection was re-established. There was no connection from %s (%s)",
			r.lastConnectionDrop.Format(time.RFC3339),
			time.Now().Sub(r.lastConnectionDrop).String(),
		)
		logrus.Warn(r.printer(result))

		r.previousResult = result
	}

	if r.previousResult.IsSuccess() && !result.IsSuccess() {
		logrus.Warnf("Connection was dropped!")
		output, err := r.printer(result)
		if err != nil {
			logrus.Errorf("Could not encode the result because of: %s", err.Error())
		} else {
			logrus.Warn(output)
		}
		r.previousResult = result
		r.lastConnectionDrop = time.Now()
	}
}
