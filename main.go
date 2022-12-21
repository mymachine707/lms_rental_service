package main

import (
	"fmt"
	"lms/config"
	"lms/protogen/rental_service"
	"lms/services/rental"
	"lms/storage"
	"lms/storage/postgres"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf := config.Load()
	AUTH := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.PostgresHost,
		conf.PostgresPort,
		conf.PostgresUser,
		conf.PostgresPassword,
		conf.PostgresDatabase,
	)
	var stg storage.StorageI
	stg, err := postgres.InitDB(AUTH)
	if err != nil {
		panic(err)
	}

	log.Printf("\ngRPC server running port: %s with tcp protocol!\n", conf.GRPCPort)

	listener, err := net.Listen("tcp", conf.GRPCPort)

	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	RentalService := rental.NewRentalService(stg)
	rental_service.RegisterRentalServiceServer(s, RentalService)

	reflection.Register(s)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
