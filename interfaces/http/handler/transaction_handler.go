package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kiosanim/pismo-code-assessment/application/transaction/dto"
	"github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/transaction"
	"net/http"
	"strconv"
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

// GetTransactionID godoc
// @Summary      Get transaction by ID
// @Description  Returns a transaction by ID
// @Tags         Transactions
// @Param        id   path	int  true  "Transaction ID"
// @Produce      json
// @Success      200  {object}  dto.FindTransactionByIdResponse
// @Failure      404  {object}  map[string]string
// @Router       /transactions/{id} [get]
func (h *TransactionHandler) GetTransactionByID(c *gin.Context) {
	transactionId := c.Param("transaction_id")
	transactionId64, err := strconv.ParseInt(transactionId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.InvalidParametersError.Error()})
		return
	}
	res, err := h.service.FindByID(c.Request.Context(), dto.FindTransactionByIdRequest{TransactionID: transactionId64})
	if res == nil && err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
