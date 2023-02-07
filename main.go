package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/TranQuocToan1996/backendMaster/api"
	db "github.com/TranQuocToan1996/backendMaster/db/sqlc"
	"github.com/TranQuocToan1996/backendMaster/gapi"
	"github.com/TranQuocToan1996/backendMaster/pb"
	"github.com/TranQuocToan1996/backendMaster/util"
	"github.com/golang-migrate/migrate/v4"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/TranQuocToan1996/backendMaster/doc/statik"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	development = "development"
	production  = "production"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	if config.Environment == development {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Print(config)

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	migrateDatabase(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)
	go gatewayServer(config, store)
	gRPCServer(config, store)
}

func migrateDatabase(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msg(err.Error())
	}

	log.Info().Msg("migration database successfully")

}

func gRPCServer(config util.Config, store db.Store) {
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	gRPCserver := grpc.NewServer(grpcLogger)
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	pb.RegisterSimpleBankServer(gRPCserver, server)
	reflection.Register(gRPCserver)

	lis, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	log.Info().Msgf("start gRPC at %v", lis.Addr())
	err = gRPCserver.Serve(lis)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func gatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg(err.Error())
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
		log.Fatal().Msg(err.Error())
	}

	statikFs, err := fs.New()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	swaggerHandle := http.StripPrefix("/swagger/", http.FileServer(statikFs))

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	mux.Handle("/swagger/", swaggerHandle)

	lis, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	log.Info().Msgf("start HTTP gateway at %v", lis.Addr().String())
	logHandler := gapi.HttpLogger(mux)
	err = http.Serve(lis, logHandler)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func ginServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	err = server.Start()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
