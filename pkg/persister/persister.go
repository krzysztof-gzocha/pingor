package persister

import (
	"github.com/krzysztof-gzocha/pingor/pkg/persister/record"
)

// Persister should be implemented by any service capable of persisting the result
type Persister interface {
	Persist(result record.Record) error
}
