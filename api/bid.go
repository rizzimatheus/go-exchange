package api

import (
	"database/sql"
	"go-exchange/util"
	"net/http"

	db "go-exchange/db/sqlc"

	"github.com/gin-gonic/gin"
)

// POST http://localhost:8080/bids
type bidRequest struct {
	Pair          string `json:"pair"`
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Price         int64  `json:"price"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createBid(ctx *gin.Context) {
	var req bidRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	c1, c2 := util.CurrenciesFromPair(req.Pair)

	_, valid := server.validAccount(ctx, req.FromAccountID, c2)
	if !valid {
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, c1)
	if !valid {
		return
	}

	arg := db.CreateBidParams{
		Pair:          req.Pair,
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Price:         req.Price,
		Amount:        req.Amount,
		Status:        "active",
	}

	result, err := server.store.CreateBid(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GET http://localhost:8080/bids/1
type getBidRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getBid(ctx *gin.Context) {
	var req getBidRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	bid, err := server.store.GetBid(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bid)
}

// GET http://localhost:8080/bids/?from_account_id=8&to_account_id=9&page_id=1&page_size=5
type listBidRequest struct {
	FromAccountID int64 `form:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `form:"to_account_id" binding:"required,min=1"`
	PageID        int32 `form:"page_id" binding:"required,min=1"`
	PageSize      int32 `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) listBids(ctx *gin.Context) {
	var req listBidRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListBidsParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Limit:         req.PageSize,
		Offset:        (req.PageID - 1) * req.PageSize,
	}

	bids, err := server.store.ListBids(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bids)
}

// PUT http://localhost:8080/bids
type updateBidRequest struct {
	ID     int64  `json:"id" binding:"required,min=1"`
	Status string `json:"status" binding:"required,eq=canceled"`
}

func (server *Server) updateBid(ctx *gin.Context) {
	var req updateBidRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateBidParams{
		ID:     req.ID,
		Status: req.Status,
	}

	bid, err := server.store.UpdateBid(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, bid)
}
