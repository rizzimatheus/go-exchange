package api

import (
	"database/sql"
	db "go-exchange/db/sqlc"
	"go-exchange/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

const(
	groupSize = int32(10)
)

func (server *Server) matchAsks(ctx *gin.Context, req bidRequest) (string, int64, error) {
	group := int32(1)
	amount := req.Amount
	var tradedAmount int64
	bidStatus := util.ACTIVE

out:
	for {
		arg := db.ListTradableAsksParams{
			Pair:   req.Pair,
			Price:  req.Price,
			Limit:  groupSize,
			Offset: (group - 1) * groupSize,
		}
		asks, err := server.store.ListTradableAsks(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return bidStatus, tradedAmount, err
		}

		if len(asks) == 0 {
			break
		}

		for _, ask := range asks {
			askAmountRemaining := ask.RemainingAmount
			askStatus := util.ACTIVE

			if ask.RemainingAmount < amount {
				amount = ask.RemainingAmount
			}
			arg := db.TradeTxParams{
				FirstFromAccountID: ask.FromAccountID,
				FirstToAccountID:   req.ToAccountID,
				FirstAmount:        amount,

				SecondFromAccountID: req.FromAccountID,
				SecondToAccountID:   ask.ToAccountID,
				SecondAmount:        ask.Price * amount,
			}

			_, err := server.store.TradeTx(ctx, arg)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return bidStatus, tradedAmount, err
			}

			//update ask
			askAmountRemaining -= amount
			if askAmountRemaining == 0 {
				askStatus = util.COMPLETED
			}
			aksArg := db.UpdateAskParams{
				ID:     ask.ID,
				Status: askStatus,
				RemainingAmount: askAmountRemaining,
			}

			_, err = server.store.UpdateAsk(ctx, aksArg)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return bidStatus, tradedAmount, err
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return bidStatus, tradedAmount, err
			}

			tradedAmount += amount
			amount = req.Amount - tradedAmount
			if amount == 0 {
				bidStatus = util.COMPLETED
				tradedAmount = 0
				break out
			}
		}

		group++
	}

	return bidStatus, tradedAmount, nil
}

func (server *Server) matchBids(ctx *gin.Context, req askRequest) (string, int64, error) {
	group := int32(1)
	amount := req.Amount
	var tradedAmount int64
	askStatus := util.ACTIVE

out:
	for {
		arg := db.ListTradableBidsParams{
			Pair:   req.Pair,
			Price:  req.Price,
			Limit:  groupSize,
			Offset: (group - 1) * groupSize,
		}
		bids, err := server.store.ListTradableBids(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return askStatus, tradedAmount, err
		}

		if len(bids) == 0 {
			break
		}

		for _, bid := range bids {
			bidAmountRemaining := bid.RemainingAmount
			bidStatus := util.ACTIVE

			if bid.RemainingAmount > amount {
				amount = bid.RemainingAmount
			}
			arg := db.TradeTxParams{
				FirstFromAccountID: req.FromAccountID,
				FirstToAccountID:   bid.ToAccountID,
				FirstAmount:        amount,

				SecondFromAccountID: bid.FromAccountID,
				SecondToAccountID:   req.ToAccountID,
				SecondAmount:        bid.Price * amount,
			}

			_, err := server.store.TradeTx(ctx, arg)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return askStatus, tradedAmount, err
			}

			//update bid
			bidAmountRemaining -= amount
			if bidAmountRemaining == 0 {
				bidStatus = util.COMPLETED
			}
			bidArg := db.UpdateBidParams{
				ID:     bid.ID,
				Status: bidStatus,
				RemainingAmount: bidAmountRemaining,
			}

			_, err = server.store.UpdateBid(ctx, bidArg)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return askStatus, tradedAmount, err
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return askStatus, tradedAmount, err
			}

			tradedAmount += amount
			amount = req.Amount - tradedAmount
			if amount == 0 {
				askStatus = util.COMPLETED
				tradedAmount = 0
				break out
			}
		}

		group++
	}

	return askStatus, tradedAmount, nil
}
