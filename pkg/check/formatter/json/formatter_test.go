// +build unit

package json

import (
	"testing"
	"time"

	"math"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/stretchr/testify/assert"
)

func TestJsonResultPrinter_Error(t *testing.T) {
	res, err := Formatter(invalidResult{Test: math.NaN()})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported value: NaN")
	assert.Equal(t, "", res)
}
func TestJsonResultPrinter_Success(t *testing.T) {
	res, err := Formatter(invalidResult{Test: 1})
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
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
func (r invalidResult) GetURL() string {
	return ""
}
func (r invalidResult) GetSubResults() []result.ResultInterface {
	return make([]result.ResultInterface, 0)
}
