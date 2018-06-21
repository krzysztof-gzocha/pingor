// +build unit

package subscriber

import (
	"testing"

	"time"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/stretchr/testify/assert"
)

func TestLogConnectionCheckResult_BadResult(t *testing.T) {
	assert.NotPanics(t, func() {
		LogConnectionCheckResult(struct{}{})
	})
}

func TestLogConnectionCheckResult_EmptyResult(t *testing.T) {
	assert.NotPanics(t, func() {
		LogConnectionCheckResult(result.Result{})
	})
}

func TestLogConnectionCheckResult_FullResult(t *testing.T) {
	assert.NotPanics(t, func() {
		LogConnectionCheckResult(result.Result{
			Time:        time.Second,
			SuccessRate: 75,
		})
	})
}
