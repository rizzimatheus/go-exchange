package api

import (
	"database/sql"
	"errors"
	db "go-exchange/db/sqlc"
	"go-exchange/token"
	"go-exchange/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST http://localhost:8080/asks
type askRequest struct {
	Pair          string `json:"pair" binding:"required,pair"`
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Price         int64  `json:"price" binding:"required,gt=0"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createAsk(ctx *gin.Context) {
	var req askRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c1, c2 := util.CurrenciesFromPair(req.Pair)

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, c1)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	toAccount, valid := server.validAccount(ctx, req.ToAccountID, c2)
	if !valid {
		return
	}

	if toAccount.Owner != authPayload.Username {
		err := errors.New("to account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	askStatus, tradedAmount, err := server.matchBids(ctx, req)
	if err != nil {
		return
	}

	arg := db.CreateAskParams{
		Pair:          req.Pair,
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Price:         req.Price,
		InitialAmount:   req.Amount,
		RemainingAmount: req.Amount - tradedAmount,
		Status:        askStatus,
	}

	result, err := server.store.CreateAsk(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GET http://localhost:8080/asks/1
type getAskRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAsk(ctx *gin.Context) {
	var req getAskRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ask, err := server.store.GetAsk(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.verifyAccountOwner(ctx, ask.FromAccountID)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, ask)
}

// GET http://localhost:8080/asks/?from_account_id=8&to_account_id=9&page_id=1&page_size=5
type listAskRequest struct {
	FromAccountID int64 `form:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `form:"to_account_id" binding:"required,min=1"`
	PageID        int32 `form:"page_id" binding:"required,min=1"`
	PageSize      int32 `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) listAsks(ctx *gin.Context) {
	var req listAskRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.verifyAccountOwner(ctx, req.FromAccountID)
	if err != nil {
		return
	}

	_, err = server.verifyAccountOwner(ctx, req.ToAccountID)
	if err != nil {
		return
	}

	arg := db.ListAsksParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Limit:         req.PageSize,
		Offset:        (req.PageID - 1) * req.PageSize,
	}

	asks, err := server.store.ListAsks(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, asks)
}

// GET http://localhost:8080/allasks/?pair=BTC_USDT&page_id=1&page_size=5
type listAllAskRequest struct {
	Pair     string `form:"pair" binding:"required,pair"`
	PageID   int32  `form:"page_id" binding:"required,min=1"`
	PageSize int32  `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) listAllAsks(ctx *gin.Context) {
	var req listAllAskRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAllAsksParams{
		Pair:   req.Pair,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	asks, err := server.store.ListAllAsks(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, asks)
}

// PUT http://localhost:8080/asks
type updateAskRequest struct {
	ID     int64  `json:"id" binding:"required,min=1"`
	Status string `json:"status" binding:"required"`
	Amount int64  `json:"amount" binding:"required,gt=0"`
}

func (server *Server) updateAsk(ctx *gin.Context) {
	var req updateAskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	a, err := server.store.GetAsk(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.verifyAccountOwner(ctx, a.FromAccountID)
	if err != nil {
		return
	}

	arg := db.UpdateAskParams{
		ID:     req.ID,
		Status: req.Status,
	}

	ask, err := server.store.UpdateAsk(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ask)
}
