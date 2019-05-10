package event

import "sync"

// Subscriber is function declaration that is able to receive an event from dispatcher
type Subscriber func(arg interface{})

// SyncedDispatcher is struct responsible for holding information about all the subscribers and inform them when proper
// event is dispatched
type SyncedDispatcher struct {
	sync.Mutex
	subscribers map[string][]Subscriber
}

// Dispatcher should be implemented by any dispatcher in the system
type Dispatcher interface {
	AttachSubscriber(eventName string, subscriber Subscriber) Dispatcher
	Dispatch(eventName string, arg interface{})
}

// NewDispatcher will return a pointer to SyncedDispatcher struct
func NewDispatcher() *SyncedDispatcher {
	return &SyncedDispatcher{
		subscribers: make(map[string][]Subscriber),
	}
}

// AttachSubscriber will add new subscriber for specific event.
func (d *SyncedDispatcher) AttachSubscriber(eventName string, subscriber Subscriber) Dispatcher {
	d.Lock()
	defer d.Unlock()

	d.subscribers[eventName] = append(d.subscribers[eventName], subscriber)

	return d
}

// Dispatch should be used to inform the dispatcher about any event. It will look for attached subscribers and run them.
// Each subscriber will run in separate go-routine
func (d *SyncedDispatcher) Dispatch(eventName string, arg interface{}) {
	d.Lock()
	defer d.Unlock()

	eventSubscribers, exists := d.subscribers[eventName]
	if !exists {
		return
	}

	for _, subscriber := range eventSubscribers {
		go subscriber(arg)
	}
}
