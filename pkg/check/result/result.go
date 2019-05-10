package result

import (
	"encoding/json"
	"time"
)

// Result should be implemented by all the results objects
type Result interface {
	IsSuccess() bool
	GetSuccessRate() float32
	GetTime() time.Duration
	GetMessage() string
	GetSubResults() []Result
	GetURL() string
}

// DefaultResult is main data transfer object used to store all the results from checkers
type DefaultResult struct {
	Success     bool          `json:"success,omitempty"`
	SuccessRate float32       `json:"success_rate"`
	Time        time.Duration `json:"time"`
	Message     string        `json:"message,omitempty"`
	SubResults  []Result      `json:"sub_results,omitempty"`
	URL         string        `json:"url,omitempty"`
}

// MarshalJSON will encode DefaultResult into JSON with time parsed into string
func (r DefaultResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Success     bool     `json:"success,omitempty"`
		SuccessRate float32  `json:"success_rate"`
		Time        string   `json:"time"`
		Message     string   `json:"message,omitempty"`
		SubResults  []Result `json:"sub_results,omitempty"`
	}{
		Success:     r.Success,
		SuccessRate: r.SuccessRate,
		Time:        r.Time.String(),
		Message:     r.Message,
		SubResults:  r.SubResults,
	})
}

// IsSuccess will return true if connection check was successful
func (r DefaultResult) IsSuccess() bool {
	return r.Success
}

// GetSuccessRate will return ratio of success check to all the checks
func (r DefaultResult) GetSuccessRate() float32 {
	return r.SuccessRate
}

// GetTime will return time taken in all the checks
func (r DefaultResult) GetTime() time.Duration {
	return r.Time
}

// GetMessage is useful to create human-readable reports
func (r DefaultResult) GetMessage() string {
	return r.Message
}

// GetSubResults will return all (if any) sub-results that were combined into this object.
// Might be useful for human-readable reports
func (r DefaultResult) GetSubResults() []Result {
	return r.SubResults
}

// GetURL will return testing URL
func (r DefaultResult) GetURL() string {
	return r.URL
}
