package service

type Report interface {
	MonthlyReport(date string) (string, error)
}
