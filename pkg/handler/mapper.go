package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) AddMoney(c *gin.Context) {
	var input addingBalanceDTO
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.service.AddMoney(input.UserID, input.MoneyAmount)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"balance": result,
	})
}

func (h *Handler) ReserveMoney(c *gin.Context) {
	var input reservingMoneyDTO
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	newReservedBalance, err := h.service.Reserve(input.UserID, input.ServiceID, input.OrderID, input.OrderCost)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"amount_reserved":      input.OrderCost,
		"new_reserved_balance": newReservedBalance,
	})
}

func (h *Handler) DealSuccess(c *gin.Context) {
	var input dealDTO
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.service.DealSuccess(input.IdUser, input.IdService, input.IdOrder, input.OrderCost)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"deal_status": true,
	})
}

func (h *Handler) GetBalance(c *gin.Context) {
	var input userIdDTO
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	balance, err := h.service.GetBalance(input.UserID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"balance": balance,
	})
}
func (h *Handler) MonthReport(c *gin.Context) {
	var input MonthReportDTO
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	csvReportPath, err := h.service.MonthlyReport(input.Date)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"report_url": csvReportPath,
	})
}
func (h *Handler) CancelOrder(c *gin.Context) {
	var input dealDTO
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.service.CancelOrder(input.IdUser, input.IdService, input.IdOrder, input.OrderCost)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"deal_canceled": true,
	})
}
func (h *Handler) UserHistory(c *gin.Context) {
	var userID userIdDTO
	if err := c.BindJSON(&userID); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	resultSet, err := h.service.UserHistory(userID.UserID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(200, resultSet)
}
