package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/TranQuocToan1996/backendMaster/api"
	db "github.com/TranQuocToan1996/backendMaster/db/sqlc"
	"github.com/TranQuocToan1996/backendMaster/gapi"
	"github.com/TranQuocToan1996/backendMaster/pb"
	"github.com/TranQuocToan1996/backendMaster/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/TranQuocToan1996/backendMaster/doc/statik"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(config)

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	go gatewayServer(config, store)
	gRPCServer(config, store)
}

func gRPCServer(config util.Config, store db.Store) {
	gRPCserver := grpc.NewServer()
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}
	pb.RegisterSimpleBankServer(gRPCserver, server)
	reflection.Register(gRPCserver)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("start gRPC at %v", lis.Addr())
	err = gRPCserver.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

func gatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	statikFs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	swaggerHandle := http.StripPrefix("/swagger/", http.FileServer(statikFs))

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	mux.Handle("/swagger/", swaggerHandle)

	lis, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("start HTTP gateway at %v", lis.Addr())
	err = http.Serve(lis, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func ginServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
