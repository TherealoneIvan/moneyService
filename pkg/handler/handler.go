package handler

import (
	"github.com/gin-gonic/gin"
	"moneyService/service"
)

type Handler struct {
	service *service.Service
}

// вызов от хендлера->логика -> передача данных к бд
func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/getBalance", h.GetBalance)
	router.GET("/addMoney", h.AddMoney)
	router.GET("/reserveMoney", h.ReserveMoney)
	router.GET("/approveDeal", h.DealSuccess)
	router.GET("/getMonthReport", h.MonthReport)
	router.GET("/cancelOrder", h.CancelOrder)
	router.GET("/getHistory", h.UserHistory)
	return router
}
