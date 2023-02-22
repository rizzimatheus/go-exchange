package api

import (
	"database/sql"
	db "go-exchange/db/sqlc"
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

	_, valid := server.validAccount(ctx, req.FromAccountID, c1)
	if !valid {
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountID, c2)
	if !valid {
		return
	}

	arg := db.CreateAskParams{
		Pair:          req.Pair,
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Price:         req.Price,
		Amount:        req.Amount,
		Status:        util.ACTIVE,
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

// PUT http://localhost:8080/asks
type updateAskRequest struct {
	ID     int64  `json:"id" binding:"required,min=1"`
	Status string `json:"status" binding:"required,eq=canceled"`
}

func (server *Server) updateAsk(ctx *gin.Context) {
	var req updateAskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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
