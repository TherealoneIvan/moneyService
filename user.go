package moneyService

import "time"

type Balance struct {
	Id              int `json:"-" db:"id"`
	Balance         int
	ReservedBalance float32
}

type UserOrders struct {
	UserId         int
	OrderID        int
	ServiceId      int
	ReservedAmount float32
}
type SucsessOrders struct {
	UserId         int
	OrderID        int
	ServiceId      int
	ReservedAmount float32
	Date           time.Time
}
