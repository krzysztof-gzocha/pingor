package mock

import (
	"github.com/krzysztof-gzocha/pingor/pkg/persister/record"
	"github.com/stretchr/testify/mock"
)

type Persister struct {
	mock.Mock
}

func (m *Persister) Persist(result record.Record) error {
	return m.Called(result).Error(0)
}
