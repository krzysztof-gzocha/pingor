package check

import (
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
	Success     bool `json:"omitempty"`
	SuccessRate float32
	Time        time.Duration
	Message     string            `json:"omitempty"`
	SubResults  []ResultInterface `json:"omitempty"`
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
