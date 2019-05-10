package result

import "time"

// MeasuredAtResult is expanding default Result by adding GetMeasuredAt method
type MeasuredAtResult interface {
	Result
	GetMeasuredAt() time.Time
}

// DefaultMeasuredAtResult is kind of wrapper for regular result with MeasuredAt field
type DefaultMeasuredAtResult struct {
	Result     Result
	MeasuredAt time.Time `json:"measured_at"`
}

// IsSuccess will return TRUE if measurement was successful
func (t DefaultMeasuredAtResult) IsSuccess() bool {
	return t.Result.IsSuccess()
}

// GetSuccessRate will return success rate >=0 and <=1 as float32
func (t DefaultMeasuredAtResult) GetSuccessRate() float32 {
	return t.Result.GetSuccessRate()
}

// GetTime will return time it took to measure. Could be useful to compare performances
func (t DefaultMeasuredAtResult) GetTime() time.Duration {
	return t.Result.GetTime()
}

// GetMessage will return additional message added to the result
func (t DefaultMeasuredAtResult) GetMessage() string {
	return t.Result.GetMessage()
}

// GetSubResults will return array of Result so it can be used as a group
func (t DefaultMeasuredAtResult) GetSubResults() []Result {
	return t.Result.GetSubResults()
}

// GetMeasuredAt will return a time where measurement was done
func (t DefaultMeasuredAtResult) GetMeasuredAt() time.Time {
	return t.MeasuredAt
}

// GetURL will return testing URL
func (t DefaultMeasuredAtResult) GetURL() string {
	return t.Result.GetURL()
}
