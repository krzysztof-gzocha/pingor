package persister

import (
	"github.com/krzysztof-gzocha/pingor/pkg/persister/record"
)

// PersisterInterface should be implemented by any service capable of persisting the result
type PersisterInterface interface {
	Persist(result record.Record) error
}
