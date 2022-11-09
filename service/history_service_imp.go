package service

import (
	"moneyService"
	"moneyService/pkg/repository"
)

type HistoryService struct {
	repository repository.HistoryOperation
}

func NewHistoryService(repository repository.HistoryOperation) *HistoryService {
	return &HistoryService{repository: repository}
}
func (h *HistoryService) UserHistory(idUser int) (map[int]moneyService.HistoryReport, error) {
	historyReport, err := h.repository.UserHistory(idUser)
	if err != nil {
		return nil, err
	}
	return historyReport, err
}
