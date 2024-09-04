package methods

import (
	"api-gateway/proto/genbooking"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CoonectBooking() genbooking.BookingServiceClient {
	conn, err := grpc.NewClient("booking_service:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connect user microservice...", err)
	}

	client := genbooking.NewBookingServiceClient(conn)
	return client
}
