package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/legobrokkori/go-kubernetes-grpc_practice/api"
	db "github.com/legobrokkori/go-kubernetes-grpc_practice/db/sqlc"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/gapi"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/pb"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configiration")
	}

	conn, err := sql.Open(config.DbDriver, config.DbServer)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store, _ := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func runGinServer(config util.Config, store db.Store) {

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
