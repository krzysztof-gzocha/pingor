package subscriber

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/pkg/check/printer"
	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
)

// Reconnection subscriber is responsible to check if connection was re-established. If so it will create proper log about it.
type Reconnection struct {
	previousResult     result.ResultInterface
	lastConnectionDrop time.Time
	printer            printer.PrinterFunc
}

// NewReconnectionSubscriber will return a pointer to Reconnection
func NewReconnectionSubscriber(printer printer.PrinterFunc) *Reconnection {
	return &Reconnection{printer: printer}
}

// NotifyAboutReconnection is subscriber method that will read the result from provided argument and interpret them.
// In case of when last result had errors and current is clear it will log this information alongside with time details.
func (r *Reconnection) NotifyAboutReconnection(arg interface{}) {
	res, ok := arg.(result.ResultInterface)
	if !ok {
		return
	}

	if r.previousResult == nil {
		r.prepareFirstPreviousResult(res)
	}

	if !r.previousResult.IsSuccess() && res.IsSuccess() {
		logrus.Warnf(
			"Connection was re-established. There was no connection from %s (%s)",
			r.lastConnectionDrop.Format(time.RFC3339),
			time.Now().Sub(r.lastConnectionDrop).String(),
		)
		r.printResult(res)
		r.previousResult = res
	}

	if r.previousResult.IsSuccess() && !res.IsSuccess() {
		logrus.Warnf("Connection was dropped!")
		r.printResult(res)
		r.previousResult = res
		r.lastConnectionDrop = time.Now()
	}
}

func (r *Reconnection) printResult(result result.ResultInterface) {
	output, err := r.printer(result)
	if err != nil {
		logrus.Errorf("Could not encode the result because of: %s", err.Error())
	} else {
		logrus.Warn(output)
	}
}

func (r *Reconnection) prepareFirstPreviousResult(result result.ResultInterface) {
	r.previousResult = result
	if !r.previousResult.IsSuccess() {
		r.lastConnectionDrop = time.Now()
	}
}
