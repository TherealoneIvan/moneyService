package service

import "moneyService/pkg/repository"

type Service struct {
	Balance
	History
	Report
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Balance: NewBalanceService(repos),
		History: NewHistoryService(repos),
		Report:  NewReportService(repos),
	}
}
