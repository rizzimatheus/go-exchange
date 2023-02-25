package gapi

import (
	"fmt"
	db "go-exchange/db/sqlc"
	"go-exchange/pb"
	"go-exchange/token"
	"go-exchange/util"
)

// Server serves gRPC requests for the exchange service.
type Server struct {
	pb.UnimplementedExchangeServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
	}

	return server, nil
}
