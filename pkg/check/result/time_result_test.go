// +build unit

package result

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeResult_IsSuccess(t *testing.T) {
	res := DefaultMeasuredAtResult{Result: DefaultResult{Success: true}}

	assert.True(t, res.IsSuccess())
}

func TestTimeResult_GetSuccessRate(t *testing.T) {
	res := DefaultMeasuredAtResult{Result: DefaultResult{SuccessRate: 0.3}}

	assert.Equal(t, float32(0.3), res.GetSuccessRate())
}

func TestTimeResult_GetTime(t *testing.T) {
	res := DefaultMeasuredAtResult{Result: DefaultResult{Time: time.Second}}

	assert.Equal(t, time.Second, res.GetTime())
}

func TestTimeResult_GetMessage(t *testing.T) {
	res := DefaultMeasuredAtResult{Result: DefaultResult{Message: "msg"}}

	assert.Equal(t, "msg", res.GetMessage())
}

func TestTimeResult_GetSubResults(t *testing.T) {
	subRes := []Result{
		DefaultResult{Success: true},
		DefaultResult{Success: true},
	}
	res := DefaultMeasuredAtResult{Result: DefaultResult{SubResults: subRes}}

	assert.Equal(t, subRes, res.GetSubResults())
}

func TestTimeResult_GetURL(t *testing.T) {
	res := DefaultMeasuredAtResult{Result: DefaultResult{URL: "test"}}
	assert.Equal(t, res.GetURL(), "test")
}

func TestTimeResult_GetMeasuredAt(t *testing.T) {
	now := time.Now()
	res := DefaultMeasuredAtResult{MeasuredAt: now}

	assert.Equal(t, now, res.GetMeasuredAt())
}
