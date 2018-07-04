package reconnection

import (
	"fmt"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/pkg/check/formatter"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber"
)

// Logger contains required things to log information about reconnection to the CLI
type Logger struct {
	pr formatter.Func
}

// NewLogger will return Logger struct
func NewLogger(pr formatter.Func) *Logger {
	return &Logger{pr: pr}
}

// LogReconnection will use logrus to log the information about reconnection
func (l Logger) LogReconnection(args interface{}) {
	event, ok := args.(subscriber.ReconnectionEvent)
	if !ok {
		return
	}

	res, err := l.pr(event.CurrentResult)
	if err != nil {
		res = fmt.Sprintf("Error during formatting: %s", err.Error())
	}

	logrus.
		WithField("lastSuccessTime", event.LastSuccess.GetMeasuredAt().Format(time.RFC3339)).
		WithField("firstConnectionDrop", event.FirstConnectionDrop.GetMeasuredAt().Format(time.RFC3339)).
		WithField("lastConnectionDrop", event.LastConnectionDrop.GetMeasuredAt().Format(time.RFC3339)).
		WithField("current", event.CurrentResult.GetMeasuredAt().Format(time.RFC3339)).
		WithField("lastSuccessRate", event.LastSuccess.GetSuccessRate()).
		WithField("currentSuccessRate", event.CurrentResult.GetSuccessRate()).
		WithField("disconnectionTime", event.DisconnectionDuration()).
		Infof("Connection was restored! Current result: %s", res)
}
