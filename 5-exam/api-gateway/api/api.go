package api

import (
	"api-gateway/methods"
	"api-gateway/proto/genbooking"
	"api-gateway/proto/genhotel"
	"api-gateway/proto/genuser"
	"context"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	user    genuser.UserServiceClient
	hotel   genhotel.HotelServiceClient
	booking genbooking.BookingServiceClient
}

func Conn() *Server {
	log.Println("Connecting to services...")
	user := methods.ConnectUser()
	hotel := methods.ConnectHotel()
	booking := methods.CoonectBooking()
	return &Server{user: user, hotel: hotel, booking: booking}
}

func RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(2, 4)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Register a new user in the system
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      genuser.RegisterReq  true  "User registration request"
// @Success      200      {object}  genuser.RegisterResp
// @Failure      400      {object}  gin.H{"error": "Invalid request"}
// @Failure      500      {object}  gin.H{"error": "Failed to register user"}
// @Router       /users/register [post]
func (s *Server) RegisterUser(c *gin.Context) {
	log.Println("Registering user...")
	var req genuser.RegisterReq

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := s.user.RegisterUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// VerifyUser godoc
// @Summary      Verify a user
// @Description  Verify user using the verification code
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      genuser.VerifyReq  true  "User verification request"
// @Success      200      {object}  genuser.VerifyResp
// @Failure      400      {object}  gin.H{"error": "Invalid request"}
// @Failure      500      {object}  gin.H{"error": "Error during verification"}
// @Router       /users/verify [post]
func (s *Server) VerifyUser(c *gin.Context) {
	log.Println("Verifying user...")
	var req genuser.VerifyReq

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := s.user.VerifyUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error during verification"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// LoginUser godoc
// @Summary      Login user
// @Description  Login a user and generate a session token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      genuser.LoginReq  true  "User login request"
// @Success      200      {object}  genuser.LoginResp
// @Failure      400      {object}  gin.H{"error": "Invalid request"}
// @Failure      500      {object}  gin.H{"error": "Error logging in user"}
// @Router       /users/login [post]
func (s *Server) LoginUser(c *gin.Context) {
	var req genuser.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := s.user.LoginUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging in user"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete a user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  genuser.DeleteUserResp
// @Failure      400  {object}  gin.H{"error": "Invalid user ID"}
// @Failure      500  {object}  gin.H{"error": "Failed to delete user"}
// @Router       /users/{id} [delete]
func (s *Server) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	req := &genuser.DeleteUserReq{UserId: int32(idInt)}
	resp, err := s.user.DeleteUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUser godoc
// @Summary      Get user
// @Description  Retrieve a user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  genuser.GetUserResp
// @Failure      400  {object}  gin.H{"error": "Invalid user ID"}
// @Failure      500  {object}  gin.H{"error": "Failed to retrieve user"}
// @Router       /users/{id} [get]
func (s *Server) GetUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	req := &genuser.GetUserReq{UserId: int32(idInt)}
	resp, err := s.user.GetUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetHotels godoc
// @Summary      Get hotels
// @Description  Retrieve a list of hotels
// @Tags         hotels
// @Produce      json
// @Success      200  {object}  genhotel.GetHotelsResponse
// @Failure      500  {object}  gin.H{"error": "Failed to retrieve hotels"}
// @Router       /hotels [get]
func (s *Server) GetHotels(c *gin.Context) {
	req := &genhotel.GetHotelsRequest{}
	resp, err := s.hotel.GetHotels(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hotels"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetHotelByID godoc
// @Summary      Get hotel details
// @Description  Retrieve details of a hotel by ID
// @Tags         hotels
// @Produce      json
// @Param        id   path      string  true  "Hotel ID"
// @Success      200  {object}  genhotel.GetHotelDetailsResponse
// @Failure      500  {object}  gin.H{"error": "Failed to retrieve hotel details"}
// @Router       /hotels/{id} [get]
func (s *Server) GetHotelByID(c *gin.Context) {
	id := c.Param("id")
	req := &genhotel.GetHotelDetailsRequest{HotelId: id}
	resp, err := s.hotel.GetHotelDetails(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hotel details"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateHotel godoc
// @Summary      Create a new hotel
// @Description  Create a new hotel record
// @Tags         hotels
// @Accept       json
// @Produce      json
// @Param        request  body      genhotel.HotelReq  true  "Hotel creation request"
// @Success      200      {object}  genhotel.HotelResp
// @Failure      400      {object}  gin.H{"error": "Invalid request"}
// @Failure      500      {object}  gin.H{"error": "Error creating hotel"}
// @Router       /hotels [post]
func (s *Server) CreateHotel(c *gin.Context) {
	var req genhotel.HotelReq
	if err := c.ShouldBind(&req); err != nil {
		log.Println("Error binding request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := s.hotel.AddHotel(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating hotel"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// CheckRoomAvailability godoc
// @Summary      Check room availability
// @Description  Check availability of rooms in a hotel
// @Tags         hotels
// @Produce      json
// @Param        id   path      string  true  "Hotel ID"
// @Success      200  {object}  genhotel.CheckRoomAvailabilityResponse
// @Failure      500  {object}  gin.H{"error": "Error checking room availability"}
// @Router       /hotels/{id}/check [get]
func (s *Server) CheckRoomAvailability(c *gin.Context) {
	id := c.Param("id")
	req := &genhotel.CheckRoomAvailabilityRequest{HotelId: id}
	resp, err := s.hotel.CheckRoomAvailability(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking room availability"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateBooking godoc
// @Summary      Create a new booking
// @Description  Create a new booking in the system
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        request  body      genbooking.BookingRequest  true  "Booking creation request"
// @Success      200      {object}  genbooking.BookingResponse
// @Failure      400      {object}  gin.H{"error": "Invalid request"}
// @Failure      500      {object}  gin.H{"error": "Error creating booking"}
// @Router       /booking [post]
func (s *Server) CreateBooking(c *gin.Context) {
	var req genbooking.BookingRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Println(&req)

	resp, err := s.booking.CreateBooking(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creat booking"})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetBookingByID godoc
// @Summary      Get booking details
// @Description  Retrieve booking details by ID
// @Tags         bookings
// @Produce      json
// @Param        id   path      int  true  "Booking ID"
// @Success      200  {object}  genbooking.BookingResponse
// @Failure      400  {object}  gin.H{"error": "Invalid booking ID"}
// @Failure      500  {object}  gin.H{"error": "Error retrieving booking"}
// @Router       /booking/{id} [get]
func (s *Server) GetBookingByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	req := &genbooking.BookingIdReq{BookingId: int32(idInt)}
	resp, err := s.booking.GetBooking(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving booking"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateBookingByID godoc
// @Summary      Update booking
// @Description  Update a booking by ID
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Param        id       path      int                           true  "Booking ID"
// @Param        request  body      genbooking.UpdateBookIdReq    true  "Booking update request"
// @Success      200      {object}  genbooking.BookingResponse
// @Failure      400      {object}  gin.H{"error": "Invalid request body"}
// @Failure      500      {object}  gin.H{"error": "Error updating booking"}
// @Router       /booking/{id} [put]
func (s *Server) UpdateBookingByID(c *gin.Context) {
	id := c.Param("id")
	bookingID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	var req genbooking.UpdateBookIdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	req.BookingId = int32(bookingID)
	resp, err := s.booking.UpdateBooking(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating booking"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	server := Conn()
	r.Use(RateLimiter())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/users/register", server.RegisterUser)
	r.POST("/users/verify", server.VerifyUser)
	r.POST("/users/login", server.LoginUser)
	r.GET("/users/:id", server.GetUser)
	r.DELETE("/users/:id", server.DeleteUser)
	r.POST("/hotels", server.CreateHotel)
	r.GET("/hotels", server.GetHotels)
	r.GET("/hotels/:id", server.GetHotelByID)
	r.GET("/hotels/:id/check", server.CheckRoomAvailability)
	r.POST("/booking", server.CreateBooking)
	r.GET("/booking/:id", server.GetBookingByID)
	r.PUT("/booking/:id", server.UpdateBookingByID)

	return r
}