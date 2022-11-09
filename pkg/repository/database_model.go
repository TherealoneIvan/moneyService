package repository

import (
	"github.com/jmoiron/sqlx"
	"moneyService"
	"time"
)

type BalanceOperation interface {
	AddingMoneyToDB(idUser int, moneyAmount float32) (float32, error)
	ReserveInDB(idUser, idService, idOrder int, orderCost, newBalance, newReservedBalance float32, orderDate time.Time) (float32, error)
	DealSuccessInDB(idUser, idService, idOrder int, moneyAmount float32, time time.Time) error
	BalanceFromDB(idUser int) (float32, error)
	ReservedBalance(idUser int) (float32, error)
	IsExists(idUser int) (bool, error)
	reserveInOrdersTable(idUser, idService, idOrder int, orderCost float32, orderDate time.Time) error
	reserveInBalanceTable(idUser int, newReservedBalance, newBalance float32) (float32, error)
	CancelOrder(idUser, idService, idOrder int, newReservedBalance, newBalance float32) error
	AddTransaction(idUser int, event string, amount float32, time time.Time) error
}
type ReportOperation interface {
	MonthlyReport(date moneyService.CustomDate) (map[int]float32, error)
}
type HistoryOperation interface {
	UserHistory(idUser int) (map[int]moneyService.HistoryReport, error)
}
type Repository struct {
	BalanceOperation
	ReportOperation
	HistoryOperation
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{BalanceOperation: NewAuthPostgres(db), ReportOperation: NewAuthPostgres(db), HistoryOperation: NewAuthPostgres(db)}
}
