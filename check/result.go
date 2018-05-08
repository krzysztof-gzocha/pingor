package check

import (
	"encoding/json"
	"time"
)

// ResultInterface should be implemented by all the results objects
type ResultInterface interface {
	IsSuccess() bool
	GetSuccessRate() float32
	GetTime() time.Duration
	GetMessage() string
	GetSubResults() []ResultInterface
}

// Result is main data transfer object used to store all the results from checkers
type Result struct {
	Success     bool              `json:"success,omitempty"`
	SuccessRate float32           `json:"success_rate"`
	Time        time.Duration     `json:"time"`
	Message     string            `json:"message,omitempty"`
	SubResults  []ResultInterface `json:"sub_results,omitempty"`
}

func (r Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Success     bool              `json:"success,omitempty"`
		SuccessRate float32           `json:"success_rate"`
		Time        string            `json:"time"`
		Message     string            `json:"message,omitempty"`
		SubResults  []ResultInterface `json:"sub_results,omitempty"`
	}{
		Success:     r.Success,
		SuccessRate: r.SuccessRate,
		Time:        r.Time.String(),
		Message:     r.Message,
		SubResults:  r.SubResults,
	})
}

// IsSuccess will return true if connection check was successful
func (r Result) IsSuccess() bool {
	return r.Success
}

// GetSuccessRate will return ratio of success check to all the checks
func (r Result) GetSuccessRate() float32 {
	return r.SuccessRate
}

// GetTime will return time taken in all the checks
func (r Result) GetTime() time.Duration {
	return r.Time
}

// GetMessage is useful to create human-readable reports
func (r Result) GetMessage() string {
	return r.Message
}

// SubResults will return all (if any) sub-results that were combined into this object.
// Might be useful for human-readable reports
func (r Result) GetSubResults() []ResultInterface {
	return r.SubResults
}
