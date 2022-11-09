package service

import "moneyService"

type History interface {
	UserHistory(idUser int) (map[int]moneyService.HistoryReport, error)
}
