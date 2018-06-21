package subscriber

import (
	"time"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/event"
)

const ReconnectionEventName = "reconnected"

// ReconnectionEvent will be used when reconnection was detected
type ReconnectionEvent struct {
	LastSuccess         result.TimeResultInterface `json:"last_success"`
	FirstConnectionDrop result.TimeResultInterface `json:"first_connection_drop"`
	LastConnectionDrop  result.TimeResultInterface `json:"last_connection_drop"`
	CurrentResult       result.TimeResultInterface `json:"current_result"`
}

// DisconnectionDuration will return disconnection duration
func (r ReconnectionEvent) DisconnectionDuration() time.Duration {
	if r.CurrentResult == nil || r.FirstConnectionDrop == nil {
		return time.Duration(0)
	}

	return r.CurrentResult.GetMeasuredAt().Sub(r.FirstConnectionDrop.GetMeasuredAt())
}

// Reconnection subscriber is responsible to check if connection was re-established. If so it will create proper log about it.
type Reconnection struct {
	dispatcher        event.DispatcherInterface
	previousResult    result.TimeResultInterface
	lastSuccessResult result.TimeResultInterface
	firstDropResult   result.TimeResultInterface
	lastDropResult    result.TimeResultInterface
}

// NewReconnection will return a pointer to Reconnection
func NewReconnection(
	dispatcher event.DispatcherInterface,
) *Reconnection {
	return &Reconnection{
		dispatcher: dispatcher,
	}
}

// NotifyAboutReconnection is subscriber method that will trigger an event when reconnection was detected
func (r *Reconnection) NotifyAboutReconnection(arg interface{}) {
	res, ok := arg.(result.TimeResultInterface)
	if !ok {
		return
	}

	if r.previousResult == nil {
		r.prepareFirstPreviousResult(res)
	}

	// Reconnected
	if !r.previousResult.IsSuccess() && res.IsSuccess() {
		r.lastDropResult = r.previousResult
		r.lastSuccessResult = res
		r.previousResult = res

		r.dispatcher.Dispatch(ReconnectionEventName, ReconnectionEvent{
			LastSuccess:         r.lastSuccessResult,
			FirstConnectionDrop: r.firstDropResult,
			LastConnectionDrop:  r.lastDropResult,
			CurrentResult:       res,
		})

		return
	}

	// Dropped
	if r.previousResult.IsSuccess() && !res.IsSuccess() {
		r.lastSuccessResult = r.previousResult
		r.previousResult = res
		r.firstDropResult = res

		return
	}

	// Still no connection
	if !r.previousResult.IsSuccess() && !res.IsSuccess() {
		r.lastDropResult = res
		r.previousResult = res

		return
	}

	// Still connected
	r.lastSuccessResult = res
	r.previousResult = res
}

func (r *Reconnection) prepareFirstPreviousResult(result result.TimeResultInterface) {
	r.previousResult = result
	if !r.previousResult.IsSuccess() {
		r.firstDropResult = result
		r.lastDropResult = result
	}

	r.lastSuccessResult = result
}
