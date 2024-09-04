package models

import (
    "google.golang.org/protobuf/types/known/timestamppb"
)


type Booking struct {
    BookingID   int32                    `json:"booking_id"`
    UserID      int32                    `json:"user_id"`
    HotelID     int32                    `json:"hotel_id"`
    RoomType    string                   `json:"room_type"`
    CheckInDate *timestamppb.Timestamp   `json:"check_in_date"`
    CheckOutDate *timestamppb.Timestamp  `json:"check_out_date"`
    TotalAmount float32                  `json:"total_amount"`
    Status      string                   `json:"status"`
}





type Only struct {
    RoomType string 
}

type Room struct {
	RoomType     string
	PricePerNight float32
}

type RoomHotel struct {
	HotelId int32
	Rooms   []Room 
}