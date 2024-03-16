package gapi

import (
	"fmt"

	db "github.com/legobrokkori/go-kubernetes-grpc_practice/db/sqlc"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/pb"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/token"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/util"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimplaBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
