package service

import (
	"booking-service/proto/genbooking"
	"booking-service/repo"
	"context"
	"log"
)

type ServerService struct {
	genbooking.UnimplementedBookingServiceServer
	r  *repo.RepoStruct
	Lg *log.Logger
}

func ConnectServ(lg *log.Logger) *ServerService {
	rr := repo.ConnectRepo()

	return &ServerService{r: rr, Lg: lg}
}

func (s *ServerService) CreateBooking(ctx context.Context, req *genbooking.BookingRequest) (*genbooking.BookingResponse, error) {
	resp, err := s.r.CreateBookingg(ctx, req)
	if err != nil {
		log.Println("Error on createbooking", err)
		return nil, err
	}

	return resp, nil
}

func (s *ServerService) GetBooking(ctx context.Context, req *genbooking.BookingIdReq) (*genbooking.BookingIdResp, error) {
	resp, err := s.r.GetBookingIdd(ctx, req)

	if err != nil {
		s.Lg.Println("Xatolik bor getbookinda", err)
		return nil, err 
	}

	return resp, nil 
}

func (s *ServerService) UpdateBooking(ctx context.Context, req *genbooking.UpdateBookIdReq) (*genbooking.UpdateBookIdResp, error) {
	resp, err := s.r.UpdateBookingg(ctx, req)

	if err != nil {
		s.Lg.Println("xatolik updatebookda", err)
		return nil, err 
	}

	return resp, nil 
}

