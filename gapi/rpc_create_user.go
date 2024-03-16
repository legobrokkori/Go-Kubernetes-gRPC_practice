package gapi

import (
	"context"

	db "github.com/legobrokkori/go-kubernetes-grpc_practice/db/sqlc"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/pb"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUserName(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "userName already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to hash create user: %s", err)
	}
	res := &pb.CreateUserResponse{
		User: converUser(user),
	}
	return res, nil
}
