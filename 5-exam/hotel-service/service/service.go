package service

import (
	"context"
	hotelpb "hotel-service/protos/genhotel"
	"hotel-service/storage/postgresql"
)

type HotelServer struct {
	storage postgresql.Just
	hotelpb.UnimplementedHotelServiceServer
}

func NewHotelServer(storage postgresql.Just) *HotelServer {
	return &HotelServer{storage: storage}
}

func (s *HotelServer) GetHotels(ctx context.Context, req *hotelpb.GetHotelsRequest) (*hotelpb.GetHotelsResponse, error) {
	hotels, err := s.storage.ListHotels(ctx)
	if err != nil {
		return nil, err
	}

	return &hotelpb.GetHotelsResponse{Hotels: hotels}, nil
}

func (s *HotelServer) AddHotel(ctx context.Context, req *hotelpb.HotelReq) (*hotelpb.AddHotelResponse, error) {
	hotelID, err := s.storage.InsertHotel(ctx, req)
	if err != nil {
		return nil, err
	}

	return &hotelpb.AddHotelResponse{
		HotelId: int32(hotelID),
	}, nil
}

func (s *HotelServer) GetHotelDetails(ctx context.Context, req *hotelpb.GetHotelDetailsRequest) (*hotelpb.GetHotelDetailsResponse, error) {
	hotelID := req.GetHotelId()
	hotel, err := s.storage.GetHotelDetails(ctx, hotelID)
	if err != nil {
		return nil, err
	}

	return &hotelpb.GetHotelDetailsResponse{Hotel: hotel}, nil
}

func (s *HotelServer) CheckRoomAvailability(ctx context.Context, req *hotelpb.CheckRoomAvailabilityRequest) (*hotelpb.CheckRoomAvailabilityResponse, error) {
	hotelID := req.GetHotelId()
	roomAvailabilities, err := s.storage.CheckRoomAvailability(ctx, hotelID)
	if err != nil {
		return nil, err
	}

	return &hotelpb.CheckRoomAvailabilityResponse{RoomAvailabilities: roomAvailabilities}, nil
}
