package api

import (
	"booking-service/logs"
	"booking-service/proto/genbooking"
	"booking-service/service"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func ConnectApi() {
	lg := logs.GetLogger("./logs/logger.log")
	log.Println("Loggeeer bu", lg)
	lis, err := net.Listen("tcp", "booking_service:8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	service := service.ConnectServ(lg)
	if service == nil {
		log.Fatalf("wrong to connect server")
	}

	grpcServer := grpc.NewServer()
	genbooking.RegisterBookingServiceServer(grpcServer, service)

	reflection.Register(grpcServer)

	log.Println("Server running on :8081")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
