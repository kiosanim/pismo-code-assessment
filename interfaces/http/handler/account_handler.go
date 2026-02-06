package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kiosanim/pismo-code-assessment/application/account/dto"
	"github.com/kiosanim/pismo-code-assessment/internal/core/errors"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"github.com/kiosanim/pismo-code-assessment/internal/domains/account"
	"net/http"
	"strconv"
)

type AccountHandler struct {
	service account.Service
	log     logger.Logger
}

func NewAccountHandler(service account.Service, log logger.Logger) *AccountHandler {
	return &AccountHandler{
		service: service,
		log:     log,
	}
}

// CreateAccount godoc
// @Summary      Create an account
// @Description  Creates a new account with a valid and not used document number
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        account  body	dto.CreateAccountRequest  true  "Account Data"
// @Success      201  {object}  dto.CreateAccountResponse
// @Failure      400  {object}  map[string]string
// @Router       /accounts [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req dto.CreateAccountRequest
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

// GetAccountByID godoc
// @Summary      Get account by ID
// @Description  Returns an account by ID
// @Tags         Accounts
// @Param        id   path	int  true  "Account ID"
// @Produce      json
// @Success      200  {object}  dto.FindAccountByIdResponse
// @Failure      404  {object}  map[string]string
// @Router       /accounts/{id} [get]
func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	accountId := c.Param("account_id")
	accountId64, err := strconv.ParseInt(accountId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.InvalidParametersError.Error()})
		return
	}
	res, err := h.service.FindByID(c.Request.Context(), dto.FindAccountByIdRequest{AccountID: accountId64})
	if res == nil && err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// ListAccounts godoc
// @Summary      List accounts with pagination
// @Description  Returns a list of accounts
// @Tags         Accounts
// @Param        cursor  path     int  false  "Pagination cursor AccountID"
// @Param        limit   path     int     false  "Max number of accounts to return (default 10)"
// @Produce      json
// @Success      200  {object}  []dto.AccountDTO
// @Failure      404  {object}  map[string]string
// @Router       /accounts/list/{cursor}/{limit} [get]
func (h *AccountHandler) ListAccounts(c *gin.Context) {
	ctx := c.Request.Context()
	// Parse query params
	limitParam := c.Param("limit")
	cursorParam := c.Param("cursor")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}
	cursor, err := strconv.Atoi(cursorParam)
	if err != nil || limit <= 0 {
		limit = 10
	}
	request := dto.ListAccountsRequest{
		Cursor: int64(cursor),
		Limit:  int64(limit),
	}
	// Call service
	response, err := h.service.List(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Build response
	resp := gin.H{
		"accounts": response,
	}
	c.JSON(http.StatusOK, resp)
}
