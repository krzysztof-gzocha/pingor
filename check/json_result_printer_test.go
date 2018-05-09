// +build unit

package check

import (
	"testing"
	"time"

	"math"

	"github.com/stretchr/testify/assert"
)

func TestJsonResultPrinter_Error(t *testing.T) {
	res, err := JsonResultPrinter(invalidResult{Test: math.NaN()})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported value: NaN")
	assert.Equal(t, "", res)
}

type invalidResult struct {
	Test float64
}

func (r invalidResult) IsSuccess() bool {
	return true
}
func (r invalidResult) GetSuccessRate() float32 {
	return 1
}
func (r invalidResult) GetTime() time.Duration {
	return time.Second
}
func (r invalidResult) GetMessage() string {
	return ""
}
func (r invalidResult) GetSubResults() []ResultInterface {
	return make([]ResultInterface, 0)
}
