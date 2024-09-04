package repo

import (
	"booking-service/handlers/hotel"
	userr "booking-service/handlers/user"
	"booking-service/proto/genbooking"
	"booking-service/storage/postgres"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type RepoStruct struct {
	Db *sql.DB
}

func ConnectRepo() *RepoStruct {
	db, err := postgres.ConnectPostgres()

	if err != nil {

		log.Println("Error on connection postgres", err)
		return nil
	}

	return &RepoStruct{Db: db}

}

func (c *RepoStruct) CreateBookingg(ctx context.Context, req *genbooking.BookingRequest) (*genbooking.BookingResponse, error) {
	query := `
    INSERT INTO hotel_book (user_id, hotel_id, room_type, check_in_date, check_out_date, total_amount, status)
    VALUES ($1, $2, $3, $4, $5, $6, 'successfully booked')
    RETURNING booking_id`

	checkInDate := req.CheckInDate.AsTime()
	checkOutDate := req.CheckOutDate.AsTime()

	days := int(checkOutDate.Sub(checkInDate).Hours() / 24)
	log.Println("Days>>>>>>>>", days)
	if days < 0 {
		return nil, fmt.Errorf("check-out date must be after check-in date")
	}

	checkInMonth, checkInDay := checkInDate.Month(), checkInDate.Day()
	checkOutMonth, checkOutDay := checkOutDate.Month(), checkOutDate.Day()
	log.Printf("Check-in Date: %s %d, Check-out Date: %s %d", checkInMonth, checkInDay, checkOutMonth, checkOutDay)

	resp, err := userr.HandleUser(int(req.UserId))
	if err != nil {
		log.Println("Error handling user:", err)
		return nil, err
	}
	log.Println("User ID:", resp)

	idStr := strconv.Itoa(int(req.HotelId))
	hotel, err := hotel.HandleHotel(idStr)
	if err != nil {
		log.Println("Error handling hotel:", err)
		return nil, err
	}
	log.Println("Hotel ID:", hotel)

	roomTypeAvailable := false
	for _, room := range hotel.Rooms {
		log.Println("------------------", room.RoomType, room.RoomType == req.RoomType, req.RoomType,"----------------------------")
		if room.RoomType == req.RoomType {
			log.Println("------------------", room.RoomType, "----------------------------")
			roomTypeAvailable = true
			break
		}
	}

	if !roomTypeAvailable {
		return nil, fmt.Errorf("room type %s not available at hotel", req.RoomType)
	}

	var roomPricePerNight float32
	for _, room := range hotel.Rooms {
		if room.RoomType == req.RoomType {
			roomPricePerNight = room.PricePerNight
			break
		}
	}
	log.Println("Roomprice >>>>", roomPricePerNight)
	totalAmount := int(roomPricePerNight) * (days)
	log.Println("totalamount>>>>", totalAmount)
	var bookingID int32

	if req.TotalAmount < float32(totalAmount) {
		return &genbooking.BookingResponse{
			Status: "Your Budget not enough sorry",
		}, nil
	}

	err = c.Db.QueryRowContext(ctx, query, req.UserId, req.HotelId, req.RoomType, checkInDate, checkOutDate, totalAmount).Scan(&bookingID)
	if err != nil {
		return nil, err
	}

	return &genbooking.BookingResponse{
		BookingId:    bookingID,
		UserId:       int32(resp),
		HotelId:      hotel.HotelId,
		RoomType:     req.RoomType,
		CheckInDate:  checkInDate.Format("2006-01-02"),
		CheckOutDate: checkOutDate.Format("2006-01-02"),
		TotalAmount:  float32(totalAmount),
		Status:       "Successfully Booked",
	}, nil
}

func (c *RepoStruct) GetBookingIdd(ctx context.Context, req *genbooking.BookingIdReq) (*genbooking.BookingIdResp, error) {
	query := `
    SELECT booking_id, user_id, hotel_id, room_type, check_in_date, check_out_date, total_amount, status
    FROM hotel_book
    WHERE booking_id = $1`

	var bookingID int32
	var userID int32
	var hotelID int32
	var roomType string
	var checkInDate, checkOutDate time.Time
	var totalAmount float32
	var status string

	err := c.Db.QueryRowContext(ctx, query, req.BookingId).Scan(
		&bookingID,
		&userID,
		&hotelID,
		&roomType,
		&checkInDate,
		&checkOutDate,
		&totalAmount,
		&status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking with id %d not found", req.BookingId)
		}
		return nil, err
	}

	checkInDateStr := checkInDate.Format("2006-01-02")
	checkOutDateStr := checkOutDate.Format("2006-01-02")

	response := &genbooking.BookingIdResp{
		BookingId:    bookingID,
		UserId:       userID,
		HotelId:      hotelID,
		RoomType:     roomType,
		CheckInDate:  checkInDateStr,
		CheckOutDate: checkOutDateStr,
		TotalAmount:  totalAmount,
		Status:       status,
	}

	return response, nil
}

func (c *RepoStruct) UpdateBookingg(ctx context.Context, req *genbooking.UpdateBookIdReq) (*genbooking.UpdateBookIdResp, error) {
	query := `
    UPDATE hotel_book
    SET check_in_date = $1, check_out_date = $2, room_type = $3
    WHERE booking_id = $4
    RETURNING booking_id, user_id, hotel_id, room_type, check_in_date, check_out_date, total_amount, status`

	checkInDate := req.CheckInDate.AsTime()
	checkOutDate := req.CheckOutDate.AsTime()

	if checkOutDate.Before(checkInDate) {
		return nil, fmt.Errorf("check-out date must be after check-in date")
	}

	var bookingID, userID, hotelID int32
	var roomType string
	var totalAmount float32
	var status string
	var updatedCheckInDate, updatedCheckOutDate time.Time

	err := c.Db.QueryRowContext(ctx, query, checkInDate, checkOutDate, req.RoomType, req.BookingId).Scan(
		&bookingID,
		&userID,
		&hotelID,
		&roomType,
		&updatedCheckInDate,
		&updatedCheckOutDate,
		&totalAmount,
		&status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking with id %d not found", req.BookingId)
		}
		return nil, err
	}

	response := &genbooking.UpdateBookIdResp{
		BookingId:    bookingID,
		UserId:       userID,
		HotelId:      hotelID,
		RoomType:     roomType,
		CheckInDate:  timestamppb.New(updatedCheckInDate),
		CheckOutDate: timestamppb.New(updatedCheckOutDate),
		TotalAmount:  totalAmount,
		Status:       status,
	}

	return response, nil
}
