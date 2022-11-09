package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"moneyService"
	"time"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

const usersBalanceTable = "users_balance"
const usersReservedOrdersTable = "users_orders"
const successOrders = "success_orders"
const allEventsTable = "all_events"

func (a *AuthPostgres) AddingMoneyToDB(idUser int, moneyAmount float32) (float32, error) {
	var balance float32
	query := fmt.Sprintf("INSERT INTO %s (id , balance) values ($1 , $2) ON CONFLICT (id) DO UPDATE SET balance = $2 RETURNING balance", usersBalanceTable)
	row := a.db.QueryRow(query, idUser, moneyAmount)
	if err := row.Scan(&balance); err != nil {
		return 0, err
	}
	return balance, nil
}
func (a *AuthPostgres) IsExists(idUser int) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS (SELECT * FROM %s WHERE id = ($1))", usersBalanceTable)
	row := a.db.QueryRow(query, idUser)
	var isExsits bool
	if err := row.Scan(&isExsits); err != nil {
		return false, err
	}
	return isExsits, nil
}
func (a *AuthPostgres) BalanceFromDB(idUser int) (float32, error) {
	userIsExists, err := a.IsExists(idUser)
	if err != nil {
		return 0, err
	}
	if !userIsExists {
		err = fmt.Errorf("user id = %d undefined ", idUser)
		return 0, err
	}
	query := fmt.Sprintf("SELECT balance FROM %s WHERE id = ($1)", usersBalanceTable)
	row := a.db.QueryRow(query, idUser)
	var balance float32
	if err := row.Scan(&balance); err != nil {
		return 0, err
	}
	return balance, nil
}
func (a *AuthPostgres) reserveInOrdersTable(idUser, idService, idOrder int, orderCost float32, orderDate time.Time) error {
	queryForOrdersTable := fmt.Sprintf("INSERT INTO %s (id , service_id , order_id , order_cost , date) values ($1 , $2 , $3 , $4 , $5)", usersReservedOrdersTable)
	row := a.db.QueryRow(queryForOrdersTable, idUser, idService, idOrder, orderCost, orderDate)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}
func (a *AuthPostgres) reserveInBalanceTable(idUser int, newReservedBalance, newBalance float32) (float32, error) {
	queryForBalance := fmt.Sprintf("UPDATE %s SET balance = ($1) , reserved_balance = ($2) WHERE id = ($3) RETURNING reserved_balance ", usersBalanceTable)
	row := a.db.QueryRow(queryForBalance, newBalance, newReservedBalance, idUser)
	var reservedBalance float32
	if err := row.Scan(&reservedBalance); err != nil {
		return 0, err
	}
	return reservedBalance, nil
}
func (a *AuthPostgres) ReserveInDB(idUser, idService, idOrder int, orderCost, newReservedBalance, newBalance float32, orderDate time.Time) (float32, error) {
	err := a.reserveInOrdersTable(idUser, idService, idOrder, orderCost, orderDate)
	if err != nil {
		return 0, err
	}
	reservedBalance, err := a.reserveInBalanceTable(idUser, newReservedBalance, newBalance)
	if err != nil {
		return 0, err
	}
	return reservedBalance, nil
}
func (a *AuthPostgres) ReservedBalance(idUser int) (float32, error) {
	query := fmt.Sprintf("SELECT reserved_balance FROM %s WHERE id = ($1)", usersBalanceTable)
	row := a.db.QueryRow(query, idUser)
	var reservedBalance float32
	if err := row.Scan(&reservedBalance); err != nil {
		return 0, err
	}
	return reservedBalance, nil
}
func (a *AuthPostgres) DealSuccessInDB(idUser, idService, idOrder int, moneyAmount float32, date time.Time) error {
	query := fmt.Sprintf("INSERT INTO %s (id ,service_id ,order_id,order_cost, date) SELECT id ,service_id ,order_id,order_cost , date FROM %s WHERE order_id = ($1)", successOrders, usersReservedOrdersTable)
	row := a.db.QueryRow(query, idOrder)
	if row.Err() != nil {
		return row.Err()
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE order_id = ($1)", usersReservedOrdersTable)
	row = a.db.QueryRow(deleteQuery, idOrder)
	if row.Err() != nil {
		return row.Err()
	}
	changeReservedBalance := fmt.Sprintf("UPDATE %s SET reserved_balance = ($1) WHERE id = ($2) ", usersBalanceTable)
	row = a.db.QueryRow(changeReservedBalance, moneyAmount, idUser)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

func (a *AuthPostgres) BalanceInDB(idUser int) (int, error) {
	return 0, nil
}

func (a *AuthPostgres) CancelOrder(idUser, idService, idOrder int, newReservedBalance, newBalance float32) error {
	query := fmt.Sprintf("UPDATE %s SET balance = ($1) , reserved_balance = ($2) WHERE id = ($3)", usersBalanceTable)
	row := a.db.QueryRow(query, newBalance, newReservedBalance, idUser)
	if row.Err() != nil {
		return row.Err()
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE order_id = ($1)", usersReservedOrdersTable)
	row = a.db.QueryRow(query, idOrder)

	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

type ReportDTO struct {
	resultSet map[int]float32
}

func (a *AuthPostgres) MonthlyReport(date moneyService.CustomDate) (map[int]float32, error) {
	resultSet := make(map[int]float32)
	reportQuery := fmt.Sprintf("SELECT service_id , SUM(order_cost) FROM %s  WHERE extract(year from date)=($1) AND extract(month from date)=($2) GROUP BY service_id", successOrders)
	rows, err := a.db.Query(reportQuery, date.Year, date.Month)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int
		var serviceSumm float32
		err = rows.Scan(&id, &serviceSumm)
		if err != nil {
			return nil, err
		}
		resultSet[id] = serviceSumm
	}

	return resultSet, nil
}

func (a *AuthPostgres) AddTransaction(idUser int, event string, amount float32, time time.Time) error {
	inputQuery := fmt.Sprintf("INSERT INTO %s  VALUES ($1 , $2 , $3 , $4);", allEventsTable)
	row := a.db.QueryRow(inputQuery, idUser, event, amount, time)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}
func (a *AuthPostgres) UserHistory(idUser int) (map[int]moneyService.HistoryReport, error) {
	resultSet := make(map[int]moneyService.HistoryReport)
	inputQuery := fmt.Sprintf("SELECT * FROM %s WHERE id = ($1);", allEventsTable)
	rows, err := a.db.Query(inputQuery, idUser)
	if err != nil {
		return nil, err
	}
	var i = 1
	for rows.Next() {

		var id int
		var moneyAmount float32
		var serviceName string
		var date time.Time
		err = rows.Scan(&id, &serviceName, &moneyAmount, &date)
		if err != nil {
			return nil, err
		}
		resultSet[i] = moneyService.HistoryReport{Event: serviceName, Amount: moneyAmount, Date: date}
		i++
	}
	fmt.Println(resultSet)
	return resultSet, nil
}
