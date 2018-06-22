package result

import "time"

type TimeResultInterface interface {
	ResultInterface
	GetMeasuredAt() time.Time
}

type TimeResult struct {
	Result     ResultInterface
	MeasuredAt time.Time `json:"measured_at"`
}

func (t TimeResult) IsSuccess() bool {
	return t.Result.IsSuccess()
}
func (t TimeResult) GetSuccessRate() float32 {
	return t.Result.GetSuccessRate()
}
func (t TimeResult) GetTime() time.Duration {
	return t.Result.GetTime()
}
func (t TimeResult) GetMessage() string {
	return t.Result.GetMessage()
}
func (t TimeResult) GetSubResults() []ResultInterface {
	return t.Result.GetSubResults()
}
func (t TimeResult) GetMeasuredAt() time.Time {
	return t.MeasuredAt
}
