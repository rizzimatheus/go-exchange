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
	case "paseto":
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
	router.POST("/users/login", server.loginUser)

	router.POST("/trades", server.createTrade)
	router.GET("/trades/:id", server.getTrade)
	router.GET("/trades", server.listTrades)

	router.GET("/transfers/:id", server.getTransfer)
	router.GET("/transfers", server.listTransfers)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.PATCH("/users", server.updateUser)
	authRoutes.DELETE("/users/:username", server.deleteUser)

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.PATCH("/accounts", server.updateAccount)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)

	authRoutes.POST("/transfers", server.createTransfer)

	authRoutes.POST("/bids", server.createBid)
	authRoutes.GET("/bids/:id", server.getBid)
	authRoutes.GET("/bids", server.listBids)
	authRoutes.PATCH("/bids", server.updateBid)

	authRoutes.POST("/asks", server.createAsk)
	authRoutes.GET("/asks/:id", server.getAsk)
	authRoutes.GET("/asks", server.listAsks)
	authRoutes.PATCH("/asks", server.updateAsk)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
