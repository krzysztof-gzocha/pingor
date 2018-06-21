// +build unit

package reconnection

import (
	"errors"
	"testing"

	"github.com/krzysztof-gzocha/pingor/pkg/mock"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber"
	"github.com/stretchr/testify/assert"
	vendorMock "github.com/stretchr/testify/mock"
)

func TestPersister_PersistReconnectionEvent_BadEvent(t *testing.T) {
	assert.NotPanics(t, func() {
		persister := &mock.Persister{}
		p := NewPersister(persister)
		p.PersistReconnectionEvent(struct{}{})

		persister.AssertNumberOfCalls(t, "Persist", 0)
	})
}

func TestPersister_PersistReconnectionEvent_PersisterError(t *testing.T) {
	persister := &mock.Persister{}
	persister.
		On("Persist", vendorMock.Anything).
		Once().
		Return(errors.New("err"))
	p := NewPersister(persister)
	p.PersistReconnectionEvent(subscriber.ReconnectionEvent{})

	persister.AssertExpectations(t)
}
