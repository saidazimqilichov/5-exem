package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	hotelpb "hotel-service/protos/genhotel"

	_ "github.com/lib/pq"
)

type Just interface {
	InsertHotel(ctx context.Context, hotel *hotelpb.HotelReq) (int, error)
	ListHotels(ctx context.Context) ([]*hotelpb.HotelDetail, error)
	GetHotelDetails(ctx context.Context, hotelID string) (*hotelpb.HotelDetail, error)
	CheckRoomAvailability(ctx context.Context, hotelID string) ([]*hotelpb.RoomAvailability, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (s *PostgresStorage) InsertHotel(ctx context.Context, hotel *hotelpb.HotelReq) (int, error) {
	if hotel.Rating < 0.0 || hotel.Rating > 11.0{
		return 0, fmt.Errorf("rating must be  0.0 and 11.0")
	}
	var hotelID int
	err := s.db.QueryRowContext(ctx,
		`INSERT INTO hotels (name, location, rating, address) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING hotel_id`,
		hotel.Name, hotel.Location, hotel.Rating, hotel.Address).Scan(&hotelID)

	if err != nil {
		return 0, err
	}

	for _, room := range hotel.Rooms {
		if room.PricePerNight < 0 || room.PricePerNight > 10000000.00 {
			return 0, fmt.Errorf("price per night must be between 0.00 and 10000000.00")
		}
		_, err := s.db.ExecContext(ctx,
			`INSERT INTO rooms (hotel_id, room_type, price_per_night, availability) 
			 VALUES ($1, $2, $3, $4)`,
			hotelID, room.RoomType, room.PricePerNight, room.Availability)
		if err != nil {
			return 0, err
		}
	}

	return hotelID, nil
}

func (s *PostgresStorage) ListHotels(ctx context.Context) ([]*hotelpb.HotelDetail, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT hotel_id, name, location, rating, address FROM hotels")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var hotels []*hotelpb.HotelDetail
	for rows.Next() {
		var hotel hotelpb.HotelDetail
		if err := rows.Scan(&hotel.HotelId, &hotel.Name, &hotel.Location, &hotel.Rating, &hotel.Address); err != nil {
			return nil, err
		}
		hotels = append(hotels, &hotel)
	}

	return hotels, nil
}

func (s *PostgresStorage) GetHotelDetails(ctx context.Context, hotelID string) (*hotelpb.HotelDetail, error) {
	query := `
		SELECT hotel_id, name, location, rating, address
		FROM hotels
		WHERE hotel_id = $1
	`
	row := s.db.QueryRowContext(ctx, query, hotelID)

	var hotel hotelpb.HotelDetail
	if err := row.Scan(&hotel.HotelId, &hotel.Name, &hotel.Location, &hotel.Rating, &hotel.Address); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("hotel with ID %s not found", hotelID)
		}
		return nil, err
	}


	roomRows, err := s.db.QueryContext(ctx, "SELECT room_type, price_per_night, availability FROM rooms WHERE hotel_id = $1", hotelID)
	if err != nil {
		return nil, err
	}
	defer roomRows.Close()

	var rooms []*hotelpb.Room
	for roomRows.Next() {
		var room hotelpb.Room
		if err := roomRows.Scan(&room.RoomType, &room.PricePerNight, &room.Availability); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}

	hotel.Rooms = rooms

	return &hotel, nil
}

func (s *PostgresStorage) CheckRoomAvailability(ctx context.Context, hotelID string) ([]*hotelpb.RoomAvailability, error) {
	query := `
		SELECT room_type, COUNT(*) as available_rooms
		FROM rooms
		WHERE hotel_id = $1 AND availability = true
		GROUP BY room_type
	`
	rows, err := s.db.QueryContext(ctx, query, hotelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roomAvailabilities []*hotelpb.RoomAvailability
	for rows.Next() {
		var roomAvailability hotelpb.RoomAvailability
		if err := rows.Scan(&roomAvailability.RoomType, &roomAvailability.AvailableRooms); err != nil {
			return nil, err
		}
		roomAvailabilities = append(roomAvailabilities, &roomAvailability)
	}

	return roomAvailabilities, nil
}
