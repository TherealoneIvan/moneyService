package handler

type addingBalanceDTO struct {
	UserID      int     `json:"id" `
	MoneyAmount float32 `json:"money_amount"`
}

type reservingMoneyDTO struct {
	UserID    int     `json:"user_id,omitempty"`
	ServiceID int     `json:"service_id,omitempty"`
	OrderID   int     `json:"order_id,omitempty"`
	OrderCost float32 `json:"order_cost,omitempty"`
}
type dealDTO struct {
	IdUser    int     `json:"id,omitempty"`
	IdService int     `json:"service_id,omitempty"`
	IdOrder   int     `json:"order_id"`
	OrderCost float32 `json:"order_cost,omitempty"`
}
type userIdDTO struct {
	UserID int `json:"id"`
}
type MonthReportDTO struct {
	Date string `json:"date"`
}
