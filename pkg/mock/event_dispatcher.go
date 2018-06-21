package mock

import (
	"github.com/krzysztof-gzocha/pingor/pkg/event"
	"github.com/stretchr/testify/mock"
)

type EventDispatcher struct {
	mock.Mock
}

func (m *EventDispatcher) AttachSubscriber(eventName string, subscriber event.Subscriber) event.DispatcherInterface {
	m.Called(eventName, subscriber)

	return m
}

func (m *EventDispatcher) Dispatch(eventName string, arg interface{}) {
	m.Called(eventName, arg)
}
