package reconnection

import (
	"github.com/Sirupsen/logrus"
	"github.com/krzysztof-gzocha/pingor/pkg/persister"
	"github.com/krzysztof-gzocha/pingor/pkg/persister/record"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber"
)

// Persister is a subscriber, that will use another perister.PersisterInterface to store the info about reconnection event
type Persister struct {
	pr persister.PersisterInterface
}

// NewPersister will return new pr subscriber, which will persist the results after reconnection
func NewPersister(persister persister.PersisterInterface) *Persister {
	return &Persister{pr: persister}
}

// PersistReconnectionEvent will call pr service to store the results around the time of reconnection
func (p *Persister) PersistReconnectionEvent(arg interface{}) {
	event, ok := arg.(subscriber.ReconnectionEvent)
	if !ok {
		return
	}

	err := p.pr.Persist(record.Transform(event))
	if err != nil {
		logrus.Errorf("Could not persist last successful result due to: %s", err.Error())
	}
}
