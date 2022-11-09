package moneyService

import "time"

type HistoryReport struct {
	Event  string
	Amount float32
	Date   time.Time
}
