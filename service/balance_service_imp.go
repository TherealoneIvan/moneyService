package service

import (
	"fmt"
	"moneyService/pkg/repository"
	"time"
)

type BalanceService struct {
	repository repository.BalanceOperation
}

func NewBalanceService(repository repository.BalanceOperation) *BalanceService {
	return &BalanceService{repository: repository}
}

type addingBalanceDTO struct {
	UserId      int `json:"id" `
	MoneyAmount int `json:"money_amount"`
}

func (s *BalanceService) AddMoney(idUser int, moneyAmount float32) (float32, error) {
	isEx, _ := s.repository.IsExists(idUser)
	err := s.repository.AddTransaction(idUser, "adding money to balance", moneyAmount, time.Now())
	if err != nil {
		return 0, err
	}
	if isEx {
		currentBalance, err := s.GetBalance(idUser)
		if err != nil {
			return 0, err
		}
		newBalance := currentBalance + moneyAmount
		return s.repository.AddingMoneyToDB(idUser, newBalance)
	}
	return s.repository.AddingMoneyToDB(idUser, moneyAmount)
}
func (s *BalanceService) Reserve(idUser, idService, idOrder int, orderCost float32) (float32, error) {
	orderdDate := time.Now()
	currentBalance, err := s.GetBalance(idUser)
	if err != nil {
		return 0, err
	}
	if currentBalance-orderCost < 0 {
		err = fmt.Errorf("insufficient balance %s", time.Now())
		return 0, err
	}
	oldReservedBalance, err := s.repository.ReservedBalance(idUser)
	if err != nil {
		return 0, err
	}
	newReservedBalance, err := s.repository.ReserveInDB(idUser, idService, idOrder,
		orderCost, oldReservedBalance+orderCost, currentBalance-orderCost, orderdDate)
	if err != nil {
		return 0, err
	}
	s.repository.AddTransaction(idUser, "reserving money from balance", orderCost, time.Now())
	return newReservedBalance, nil
}
func (s *BalanceService) DealSuccess(idUser, idService, idOrder int, moneyAmount float32) error {
	orderdSuccsesTime := time.Now()
	oldReservedMoneyAmount, err := s.repository.ReservedBalance(idUser)
	if err != nil {
		return err
	}
	newReservedMoneyAmount := oldReservedMoneyAmount - moneyAmount
	fmt.Println(newReservedMoneyAmount)
	err = s.repository.DealSuccessInDB(idUser, idService, idOrder, newReservedMoneyAmount, orderdSuccsesTime)
	if err != nil {
		return err
	}
	s.repository.AddTransaction(idUser, "successful payment", moneyAmount, orderdSuccsesTime)
	return nil
}

func (s *BalanceService) GetBalance(idUser int) (float32, error) {
	balance, err := s.repository.BalanceFromDB(idUser)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
func (s *BalanceService) CancelOrder(idUser, idService, idOrder int, moneyAmount float32) error {
	previousReservedBalance, err := s.repository.ReservedBalance(idUser)
	if err != nil {
		return err
	}
	newReservedBalance := previousReservedBalance - moneyAmount
	previousBalance, err := s.GetBalance(idUser)
	if err != nil {
		return err
	}
	newBalance := previousBalance + moneyAmount
	s.repository.CancelOrder(idUser, idService, idOrder, newReservedBalance, newBalance)
	s.repository.AddTransaction(idUser, "canceling order", moneyAmount, time.Now())
	return nil
}
