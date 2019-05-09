// +build unit

package result

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResult_IsSuccess(t *testing.T) {
	res := Result{Success: true}
	assert.True(t, res.IsSuccess())
}

func TestResult_GetMessage(t *testing.T) {
	res := Result{Message: "msg"}
	assert.Equal(t, "msg", res.GetMessage())
}

func TestResult_GetSuccessRate(t *testing.T) {
	res := Result{SuccessRate: 0.6}
	assert.Equal(t, float32(0.6), res.GetSuccessRate())
}

func TestResult_GetTime(t *testing.T) {
	res := Result{Time: time.Second}
	assert.Equal(t, time.Second, res.GetTime())
}

func TestResult_GetSubResults(t *testing.T) {
	res := Result{SubResults: []ResultInterface{Result{Success: true}, Result{Success: true}}}
	assert.Len(t, res.GetSubResults(), 2)
}
func TestResult_GetURL(t *testing.T) {
	res := Result{URL: "test"}
	assert.Equal(t, res.GetURL(), "test")
}

func TestResult_MarshalJSON(t *testing.T) {
	originalResult := Result{
		Message:     "Message",
		SubResults:  make([]ResultInterface, 0),
		Success:     true,
		SuccessRate: 0.75,
		Time:        time.Second + time.Millisecond,
	}
	jsonEncoded, err := json.Marshal(originalResult)
	assert.Nil(t, err)
	decodedResult := struct {
		Success     bool              `json:"success,omitempty"`
		SuccessRate float32           `json:"success_rate"`
		Time        string            `json:"time"`
		Message     string            `json:"message,omitempty"`
		SubResults  []ResultInterface `json:"sub_results,omitempty"`
	}{}

	err = json.Unmarshal([]byte(jsonEncoded), &decodedResult)
	assert.Nil(t, err)
	assert.Equal(t, originalResult.SuccessRate, decodedResult.SuccessRate)
	assert.Equal(t, originalResult.Message, decodedResult.Message)
	assert.Equal(t, originalResult.Success, decodedResult.Success)
	assert.Equal(t, "1.001s", decodedResult.Time)
}
