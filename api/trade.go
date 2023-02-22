package api

import (
	"database/sql"
	db "go-exchange/db/sqlc"
	"go-exchange/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tradeRequest struct {
	FirstFromAccountID int64 `json:"first_from_account_id" binding:"required,min=1"`
	FirstToAccountID   int64 `json:"first_to_account_id" binding:"required,min=1"`
	FirstAmount        int64 `json:"first_amount" binding:"required,gt=0"`

	SecondFromAccountID int64 `json:"second_from_account_id" binding:"required,min=1"`
	SecondToAccountID   int64 `json:"second_to_account_id" binding:"required,min=1"`
	SecondAmount        int64 `json:"second_amount" binding:"required,gt=0"`

	Pair string `json:"pair" binding:"required,pair"`
}

func (server *Server) createTrade(ctx *gin.Context) {
	var req tradeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c1, c2 := util.CurrenciesFromPair(req.Pair)

	_, valid := server.validAccount(ctx, req.FirstFromAccountID, c1)
	if !valid {
		return
	}

	_, valid = server.validAccount(ctx, req.FirstToAccountID, c1)
	if !valid {
		return
	}

	_, valid = server.validAccount(ctx, req.SecondFromAccountID, c2)
	if !valid {
		return
	}

	_, valid = server.validAccount(ctx, req.SecondToAccountID, c2)
	if !valid {
		return
	}

	arg := db.TradeTxParams{
		FirstFromAccountID: req.FirstFromAccountID,
		FirstToAccountID:   req.FirstToAccountID,
		FirstAmount:        req.FirstAmount,

		SecondFromAccountID: req.SecondFromAccountID,
		SecondToAccountID:   req.SecondToAccountID,
		SecondAmount:        req.SecondAmount,
	}

	result, err := server.store.TradeTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GET http://localhost:8080/trades/1
type getTradeRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTrade(ctx *gin.Context) {
	var req getTradeRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	trade, err := server.store.GetTrade(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, trade)
}

// GET http://localhost:8080/trades/?from_account_id=8&to_account_id=9&page_id=1&page_size=5
type listTradeRequest struct {
	FirstFromAccountID int64 `form:"first_from_account_id" binding:"required,min=1"`
	FirstToAccountID   int64 `form:"first_to_account_id" binding:"required,min=1"`

	SecondFromAccountID int64 `form:"second_from_account_id" binding:"required,min=1"`
	SecondToAccountID   int64 `form:"second_to_account_id" binding:"required,min=1"`

	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) listTrades(ctx *gin.Context) {
	var req listTradeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListTradesParams{
		FirstFromAccountID: req.FirstFromAccountID,
		FirstToAccountID:   req.FirstToAccountID,

		SecondFromAccountID: req.SecondFromAccountID,
		SecondToAccountID:   req.SecondToAccountID,

		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	trades, err := server.store.ListTrades(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, trades)
}
