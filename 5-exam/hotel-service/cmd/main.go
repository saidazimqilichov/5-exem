package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"hotel-service/config"
	connpostgresql "hotel-service/pkg/connPostgreSql"
	hotelpb "hotel-service/protos/genhotel"
	"hotel-service/service"
	"hotel-service/storage/postgresql"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load(".")
	db, err := connpostgresql.ConnectToDB(cfg)
	log.Println("Xatolik>>>>>", err)
	if err != nil {
		log.Fatalf("wrongggg to connect to the database: %v", err)
	}
	storage := postgresql.NewPostgresStorage(db)

	hotelservice := service.NewHotelServer(storage)

	lis, err := net.Listen("tcp", ":"+cfg.ServicePort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpcInterceptor),
	)
	reflection.Register(s)

	hotelpb.RegisterHotelServiceServer(s, hotelservice)

	fmt.Printf("Service listening on port :%s\n", cfg.ServicePort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


func grpcInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	m, err := handler(ctx, req)
	if err != nil {
		log.Printf("RPC failed with error: %v", err)
		return nil, err
	}
	return m, nil
}
