package gapi

import (
	db "github.com/legobrokkori/go-kubernetes-grpc_practice/db/sqlc"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func converUser(user db.User) *pb.User {
	return &pb.User{
		UserName:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		PasswordChangeAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:        timestamppb.New(user.CreatedAt),
	}
}
