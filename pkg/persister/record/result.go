package record

import (
	"time"

	"github.com/krzysztof-gzocha/pingor/pkg/check/result"
	"github.com/krzysztof-gzocha/pingor/pkg/subscriber"
)

// Record is a struct which hold all the information about reconnection event
type Record struct {
	DeviceName       string                     `json:"device"`
	Date             time.Time                  `json:"date"`
	ConnDropDuration time.Duration              `json:"conn_drop_duration"`
	LastSuccess      result.TimeResultInterface `json:"last_success"`
	FirstConnDrop    result.TimeResultInterface `json:"first_drop"`
	CurrentResult    result.TimeResultInterface `json:"current_result"`
}

// Transform will transform ReconnectionEvent into Record
func Transform(event subscriber.ReconnectionEvent) Record {
	date := time.Now()
	if event.CurrentResult != nil {
		date = event.CurrentResult.GetMeasuredAt()
	}

	return Record{
		Date:             date,
		ConnDropDuration: event.DisconnectionDuration(),
		LastSuccess:      event.LastSuccess,
		FirstConnDrop:    event.FirstConnectionDrop,
		CurrentResult:    event.CurrentResult,
	}
}
