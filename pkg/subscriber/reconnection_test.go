// +build unit

package subscriber

import (
	"testing"

	"time"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestReconnectionEvent_DisconnectionDuration(t *testing.T) {
	now := time.Now()
	after := now.Add(time.Second)
	e := ReconnectionEvent{
		CurrentResult:       result.TimeResult{MeasuredAt: after},
		FirstConnectionDrop: result.TimeResult{MeasuredAt: now},
	}

	assert.Equal(t, time.Second, e.DisconnectionDuration())
}

func TestReconnection_NotifyAboutReconnection_BadArgument(t *testing.T) {
	r := NewReconnection(&mock.EventDispatcher{})
	assert.NotPanics(t, func() {
		r.NotifyAboutReconnection(struct{}{})
	})
}

func TestReconnection_NotifyAboutReconnection_Reconnected(t *testing.T) {
	now := time.Now()
	badResult := result.TimeResult{
		MeasuredAt: now,
		Result:     result.Result{Success: false},
	}
	goodResult := result.TimeResult{
		MeasuredAt: now,
		Result:     result.Result{Success: true},
	}
	expectedEvent := ReconnectionEvent{
		LastSuccess:         goodResult,
		CurrentResult:       goodResult,
		LastConnectionDrop:  badResult,
		FirstConnectionDrop: badResult,
	}

	eventDispatcher := &mock.EventDispatcher{}
	eventDispatcher.
		On("Dispatch", ReconnectionEventName, expectedEvent).
		Once()

	r := NewReconnection(eventDispatcher)

	r.NotifyAboutReconnection(badResult)
	eventDispatcher.AssertNumberOfCalls(t, "Dispatch", 0)

	r.NotifyAboutReconnection(goodResult)
	eventDispatcher.AssertNumberOfCalls(t, "Dispatch", 1)
	eventDispatcher.AssertExpectations(t)
}

func TestReconnection_NotifyAboutReconnection_ReconnectedSeveralChecks(t *testing.T) {
	now := time.Now()
	badResult := result.TimeResult{
		MeasuredAt: now,
		Result:     result.Result{Success: false},
	}
	secondBadResult := result.TimeResult{
		MeasuredAt: now,
		Result:     result.Result{Success: false},
	}
	goodResult := result.TimeResult{
		MeasuredAt: now,
		Result:     result.Result{Success: true},
	}
	expectedEvent := ReconnectionEvent{
		LastSuccess:         goodResult,
		CurrentResult:       goodResult,
		LastConnectionDrop:  secondBadResult,
		FirstConnectionDrop: badResult,
	}

	eventDispatcher := &mock.EventDispatcher{}
	eventDispatcher.
		On("Dispatch", ReconnectionEventName, expectedEvent).
		Once()

	r := NewReconnection(eventDispatcher)

	r.NotifyAboutReconnection(badResult)
	eventDispatcher.AssertNumberOfCalls(t, "Dispatch", 0)

	r.NotifyAboutReconnection(secondBadResult)
	eventDispatcher.AssertNumberOfCalls(t, "Dispatch", 0)

	r.NotifyAboutReconnection(goodResult)
	eventDispatcher.AssertNumberOfCalls(t, "Dispatch", 1)
	eventDispatcher.AssertExpectations(t)
}

func TestReconnection_NotifyAboutReconnection_ConnectionDropped(t *testing.T) {
	now := time.Now()
	badResult := result.TimeResult{
		MeasuredAt: now,
		Result:     result.Result{Success: false},
	}
	goodResult := result.TimeResult{
		MeasuredAt: now,
		Result:     result.Result{Success: true},
	}

	eventDispatcher := &mock.EventDispatcher{}

	r := NewReconnection(eventDispatcher)

	r.NotifyAboutReconnection(goodResult)
	eventDispatcher.AssertNumberOfCalls(t, "Dispatch", 0)

	r.NotifyAboutReconnection(badResult)
	eventDispatcher.AssertNumberOfCalls(t, "Dispatch", 0)
}

func TestReconnection_NotifyAboutReconnection_StillNoConnection(t *testing.T) {
	now := time.Now()
	badResult := result.TimeResult{
		MeasuredAt: now,
		Result:     result.Result{Success: false},
	}

	eventDispatcher := &mock.EventDispatcher{}

	r := NewReconnection(eventDispatcher)

	r.NotifyAboutReconnection(badResult)
	eventDispatcher.AssertNumberOfCalls(t, "Dispatch", 0)

	r.NotifyAboutReconnection(badResult)
	eventDispatcher.AssertNumberOfCalls(t, "Dispatch", 0)
	eventDispatcher.AssertExpectations(t)
}

func TestReconnection_NotifyAboutReconnection_StillConnected(t *testing.T) {
	now := time.Now()
	badResult := result.TimeResult{
		MeasuredAt: now,
		Result:     result.Result{Success: true},
	}

	eventDispatcher := &mock.EventDispatcher{}

	r := NewReconnection(eventDispatcher)

	r.NotifyAboutReconnection(badResult)
	eventDispatcher.AssertNumberOfCalls(t, "Dispatch", 0)

	r.NotifyAboutReconnection(badResult)
	eventDispatcher.AssertNumberOfCalls(t, "Dispatch", 0)
	eventDispatcher.AssertExpectations(t)
}
