package gapi

import (
	"context"

	db "github.com/legobrokkori/go-kubernetes-grpc_practice/db/sqlc"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/pb"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/util"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/val"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
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

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUserName(req.GetUserName()); err != nil {
		violations = append(violations, fieldViolation("user_name", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
