package result

import "time"

// TimeResultInterface is expanding default ResultInterface by adding GetMeasuredAt method
type TimeResultInterface interface {
	ResultInterface
	GetMeasuredAt() time.Time
}

// TimeResult is kind of wrapper for regular Result with MeasuredAt field
type TimeResult struct {
	Result     ResultInterface
	MeasuredAt time.Time `json:"measured_at"`
}

// IsSuccess will return TRUE if measurement was successful
func (t TimeResult) IsSuccess() bool {
	return t.Result.IsSuccess()
}

// GetSuccessRate will return success rate >=0 and <=1 as float32
func (t TimeResult) GetSuccessRate() float32 {
	return t.Result.GetSuccessRate()
}

// GetTime will return time it took to measure. Could be useful to compare performances
func (t TimeResult) GetTime() time.Duration {
	return t.Result.GetTime()
}

// GetMessage will return additional message added to the result
func (t TimeResult) GetMessage() string {
	return t.Result.GetMessage()
}

// GetSubResults will return array of ResultInterface so it can be used as a group
func (t TimeResult) GetSubResults() []ResultInterface {
	return t.Result.GetSubResults()
}

// GetMeasuredAt will return a time where measurement was done
func (t TimeResult) GetMeasuredAt() time.Time {
	return t.MeasuredAt
}
