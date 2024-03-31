package api

import (
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(c *gin.Context) {
	var req createTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	fmt.Println(arg, req)
	if !server.validateAccount(c, req.FromAccountID, req.Currency) {
		return
	}
	if !server.validateAccount(c, req.ToAccountID, req.Currency) {
		return
	}
	result, err := server.store.TransferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, result)
}

func (server *Server) validateAccount(c *gin.Context, accountId int64, currency string) bool {
	fmt.Println(accountId, currency)
	account, err := server.store.GetAccount(c, accountId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if account.Currency != currency {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
