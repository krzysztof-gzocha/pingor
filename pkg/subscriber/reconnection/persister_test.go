// +build unit

package reconnection

import (
	"errors"
	"testing"

	"github.com/krzysztof-gzocha/pingor/pkg/persister/record"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPersister_PersistReconnectionEvent_BadEvent(t *testing.T) {
	assert.NotPanics(t, func() {
		persister := &persisterMock{}
		p := NewPersister(persister)
		p.PersistReconnectionEvent(struct{}{})

		persister.AssertNumberOfCalls(t, "Persist", 0)
	})
}

func TestPersister_PersistReconnectionEvent_PersisterError(t *testing.T) {
	persister := &persisterMock{}
	persister.
		On("Persist", mock.Anything).
		Once().
		Return(errors.New("err"))
	p := NewPersister(persister)
	p.PersistReconnectionEvent(subscriber.ReconnectionEvent{})

	persister.AssertExpectations(t)
}

type persisterMock struct {
	mock.Mock
}

func (m *persisterMock) Persist(result record.Record) error {
	return m.Called(result).Error(0)
}
