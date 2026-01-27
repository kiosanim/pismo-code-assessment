package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
	"github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"net/http"
)

type TransactionHandler struct {
	service transaction.Service
	log     logger.Logger
}

func NewTransactionHandler(service transaction.Service, log logger.Logger) *TransactionHandler {
	return &TransactionHandler{
		service: service,
		log:     log,
	}
}

// CreateAccount godoc
// @Summary      Create a transaction
// @Description  Creates a new transaction with valid account id and document number
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        account  body	dto.CreateTransactionRequest  true  "Transaction Data"
// @Success      201  {object}  dto.CreateTransactionResponse
// @Failure      400  {object}  map[string]string
// @Router       /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.InvalidParametersError.Error()})
		return
	}
	res, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}
