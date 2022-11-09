package service

import (
	"encoding/csv"
	"fmt"
	"moneyService"
	"moneyService/pkg/repository"
	"os"
	"strconv"
	"strings"
)

const csvPath = "./report.csv"

type ReportService struct {
	repository repository.ReportOperation
}

func NewReportService(repository repository.ReportOperation) *ReportService {
	return &ReportService{repository: repository}
}

func stringToDateFormat(date string) (moneyService.CustomDate, error) {
	splitDate := strings.Split(date, "-")
	year, err := strconv.Atoi(splitDate[0])
	if err != nil {
		return moneyService.CustomDate{}, err
	}
	month, err := strconv.Atoi(splitDate[1])
	if err != nil {
		return moneyService.CustomDate{}, err
	}
	return moneyService.CustomDate{
		Year:  year,
		Month: month,
	}, nil
}
func (r *ReportService) MonthlyReport(date string) (string, error) {
	customDate, err := stringToDateFormat(date)
	fmt.Println(customDate.Month)
	if err != nil {
		return "", err
	}
	fmt.Println(-1)
	monthlyReport, err := r.repository.MonthlyReport(customDate)
	if err != nil {
		return "", err
	}
	file, err := os.Create(csvPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	fmt.Println(monthlyReport)
	for key, value := range monthlyReport {
		var r []string
		r = append(r, fmt.Sprintf("service_id: %v", key))
		r = append(r, fmt.Sprintf("total_money_amount: %v", value))
		err := w.Write(r)
		fmt.Println(r)
		if err != nil {
			return "", err
		}
	}
	w.Flush()
	return csvPath, nil
}
