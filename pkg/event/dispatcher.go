package event

import "sync"

// Subscriber is function declaration that is able to receive an event from dispatcher
type Subscriber func(arg interface{})

// Dispatcher is struct responsible for holding information about all the subscribers and inform them when proper
// event is dispatched
type Dispatcher struct {
	sync.Mutex
	subscribers map[string][]Subscriber
}

// DispatcherInterface should be implemented by any dispatcher in the system
type DispatcherInterface interface {
	AttachSubscriber(eventName string, subscriber Subscriber) DispatcherInterface
	Dispatch(eventName string, arg interface{})
}

// NewDispatcher will return a pointer to Dispatcher struct
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		subscribers: make(map[string][]Subscriber),
	}
}

// AttachSubscriber will add new subscriber for specific event.
func (d *Dispatcher) AttachSubscriber(eventName string, subscriber Subscriber) DispatcherInterface {
	d.Lock()
	defer d.Unlock()

	d.subscribers[eventName] = append(d.subscribers[eventName], subscriber)

	return d
}

// Dispatch should be used to inform the dispatcher about any event. It will look for attached subscribers and run them.
// Each subscriber will run in separate go-routine
func (d *Dispatcher) Dispatch(eventName string, arg interface{}) {
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
