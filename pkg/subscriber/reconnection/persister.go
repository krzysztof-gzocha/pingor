package reconnection

import (
	"github.com/krzysztof-gzocha/pingor/pkg/log"
	"github.com/krzysztof-gzocha/pingor/pkg/persister"
	"github.com/krzysztof-gzocha/pingor/pkg/persister/record"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber"
)

// Persister is a subscriber, that will use another perister.Persister to store the info about reconnection event
type Persister struct {
	logger log.Logger
	pr     persister.Persister
}

// NewPersister will return new pr subscriber, which will persist the results after reconnection
func NewPersister(logger log.Logger, persister persister.Persister) *Persister {
	return &Persister{logger: logger, pr: persister}
}

// PersistReconnectionEvent will call pr service to store the results around the time of reconnection
func (p *Persister) PersistReconnectionEvent(arg interface{}) {
	event, ok := arg.(subscriber.ReconnectionEvent)
	if !ok {
		return
	}

	err := p.pr.Persist(record.Transform(event))
	if err != nil {
		p.logger.Errorf("Could not persist last successful result due to: %s", err.Error())
	}
}
