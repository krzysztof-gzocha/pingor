package mock

import (
	"github.com/krzysztof-gzocha/pingor/pkg/event"
	"github.com/stretchr/testify/mock"
)

// EventDispatcher is just a mock
type EventDispatcher struct {
	mock.Mock
}

// AttachSubscriber is just a mock
func (m *EventDispatcher) AttachSubscriber(eventName string, subscriber event.Subscriber) event.Dispatcher {
	m.Called(eventName, subscriber)

	return m
}

// Dispatch is just a mock
func (m *EventDispatcher) Dispatch(eventName string, arg interface{}) {
	m.Called(eventName, arg)
}
