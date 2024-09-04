package methods

import (
	"api-gateway/proto/genhotel"
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ConnectHotel() genhotel.HotelServiceClient {
	conn, err := grpc.NewClient("hotel_service:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connect user microservice...", err)
	}

	client := genhotel.NewHotelServiceClient(conn)
	return client
}
