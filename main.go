package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/api"
	db "github.com/legobrokkori/go-kubernetes-grpc_practice/db/sqlc"
	_ "github.com/legobrokkori/go-kubernetes-grpc_practice/doc/statik"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/gapi"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/pb"
	"github.com/legobrokkori/go-kubernetes-grpc_practice/util"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configiration", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DbServer)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store, _ := db.NewStore(conn)
	go runGatewayServer(config, store)
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
		log.Fatal("cannot create listener", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal("cannot register handler server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik fs: ", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot create listener", err)
	}

	log.Printf("start HTTP server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server", err)
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
