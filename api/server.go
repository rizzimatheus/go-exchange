package api

import (
	"fmt"
	db "go-exchange/db/sqlc"
	"go-exchange/token"
	"go-exchange/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for exchange service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	var tokenMaker token.Maker
	var err error

	switch config.TokenType {
	case "jwt":
		tokenMaker, err = token.NewJWTMaker(config.TokenSymmetricKey)
	case "passeto":
		tokenMaker, err = token.NewPasetoMaker(config.TokenSymmetricKey)
	default:
		tokenMaker, err = token.NewPasetoMaker(config.TokenSymmetricKey)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
		v.RegisterValidation("pair", validPair)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.PATCH("/users", server.updateUser)
	router.DELETE("/users/:username", server.deleteUser)
	router.POST("/users/login", server.loginUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PATCH("/accounts", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers", server.createTransfer)
	router.GET("/transfers/:id", server.getTransfer)
	router.GET("/transfers", server.listTransfers)

	router.POST("/trades", server.createTrade)
	router.GET("/trades/:id", server.getTrade)
	router.GET("/trades", server.listTrades)

	router.POST("/bids", server.createBid)
	router.GET("/bids/:id", server.getBid)
	router.GET("/bids", server.listBids)
	router.PATCH("/bids", server.updateBid)

	router.POST("/asks", server.createAsk)
	router.GET("/asks/:id", server.getAsk)
	router.GET("/asks", server.listAsks)
	router.PATCH("/asks", server.updateAsk)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
