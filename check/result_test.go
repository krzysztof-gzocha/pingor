// +build unit

package check

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResult_MarshalJSON(t *testing.T) {
	originalResult := Result{
		Message:     "Message",
		SubResults:  make([]ResultInterface, 0),
		Success:     true,
		SuccessRate: 0.75,
		Time:        time.Second + time.Millisecond,
	}
	jsonEncoded, err := JsonResultPrinter(originalResult)
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
