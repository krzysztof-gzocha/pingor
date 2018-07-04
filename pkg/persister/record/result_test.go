// +build unit

package record

import (
	"testing"

	"time"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber"
	"github.com/stretchr/testify/assert"
)

func TestTransform_EmptyCurrentResult(t *testing.T) {
	event := subscriber.ReconnectionEvent{}

	res := Transform(event)

	assert.NotEmpty(t, res.Date)
}

func TestTransform_ExistingCurrentResult(t *testing.T) {
	now := time.Now()
	event := subscriber.ReconnectionEvent{
		CurrentResult: result.TimeResult{
			MeasuredAt: now,
		},
	}

	res := Transform(event)

	assert.NotEmpty(t, res.Date)
	assert.Equal(t, now, res.Date)
}
