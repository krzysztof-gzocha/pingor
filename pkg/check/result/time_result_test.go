// +build unit

package result

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeResult_IsSuccess(t *testing.T) {
	res := TimeResult{Result: Result{Success: true}}

	assert.True(t, res.IsSuccess())
}

func TestTimeResult_GetSuccessRate(t *testing.T) {
	res := TimeResult{Result: Result{SuccessRate: 0.3}}

	assert.Equal(t, float32(0.3), res.GetSuccessRate())
}

func TestTimeResult_GetTime(t *testing.T) {
	res := TimeResult{Result: Result{Time: time.Second}}

	assert.Equal(t, time.Second, res.GetTime())
}

func TestTimeResult_GetMessage(t *testing.T) {
	res := TimeResult{Result: Result{Message: "msg"}}

	assert.Equal(t, "msg", res.GetMessage())
}

func TestTimeResult_GetSubResults(t *testing.T) {
	subRes := []ResultInterface{
		Result{Success: true},
		Result{Success: true},
	}
	res := TimeResult{Result: Result{SubResults: subRes}}

	assert.Equal(t, subRes, res.GetSubResults())
}

func TestTimeResult_GetMeasuredAt(t *testing.T) {
	now := time.Now()
	res := TimeResult{MeasuredAt: now}

	assert.Equal(t, now, res.GetMeasuredAt())
}
