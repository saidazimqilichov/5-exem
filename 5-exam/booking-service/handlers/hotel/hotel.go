package hotel

import (
	"booking-service/models"
	"booking-service/proto/genhotel"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func HandleHotel(id string) (models.RoomHotel, error) {
	log.Println("HandleHotel function started")
	conn, err := grpc.NewClient("hotel_service:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Error connecting to gRPC server:", err)
		return models.RoomHotel{}, err
	}
	defer conn.Close()

	client := genhotel.NewHotelServiceClient(conn)

	hotelResponse, err := client.GetHotelDetails(context.Background(), &genhotel.GetHotelDetailsRequest{HotelId: id})
	if err != nil {
		log.Println("Error getting hotel details:", err)
		return models.RoomHotel{}, err
	}

	
	hotel := models.RoomHotel{
		HotelId: int32(hotelResponse.Hotel.HotelId),
		Rooms:   make([]models.Room, len(hotelResponse.Hotel.Rooms)),
	}

	for i, room := range hotelResponse.Hotel.Rooms {
		hotel.Rooms[i] = models.Room{
			RoomType:     room.RoomType,
			PricePerNight: float32(room.PricePerNight),
		}
	}

	return hotel, nil
}
