package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"net/http"
)

type TransactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}
