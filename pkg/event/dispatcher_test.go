// +build unit

package event

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewDispatcher(t *testing.T) {
	dispatcher := NewDispatcher()
	assert.Len(t, dispatcher.subscribers, 0)
}

func TestDispatcher_AttachSubscriber(t *testing.T) {
	dispatcher := NewDispatcher()
	dispatcher.AttachSubscriber("eventName", func(arg interface{}) {})
	assert.Len(t, dispatcher.subscribers, 1)
}

func TestDispatcher_Dispatch(t *testing.T) {
	dispatcher := NewDispatcher()
	subscriber := subscriberMock{}
	event := event{}
	subscriber.On("OnEvent", event).Once()
	dispatcher.AttachSubscriber("eventName", subscriber.OnEvent)
	dispatcher.Dispatch("eventName", event)
	<-time.After(time.Millisecond * 100)
	assert.True(t, subscriber.AssertExpectations(t))

	assert.NotPanics(t, func() {
		dispatcher.Dispatch("eventWithoutSubscriber", event)
	})
}

type event struct{}

type subscriberMock struct {
	mock.Mock
}

func (m subscriberMock) OnEvent(arg interface{}) {
	m.Called(arg)
}
