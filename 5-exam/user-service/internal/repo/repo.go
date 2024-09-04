package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"user-service/internal/email"
	"user-service/internal/hash"
	"user-service/internal/jwt"
	"user-service/models"
	"user-service/protos/genuser"
	postgress "user-service/storage/postgres"
	rediss "user-service/storage/redis"

	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	Db  *sql.DB
	Rdb *redis.Client
	// Kafka *kafka.Kafka
}

func Connect() *UserServer {
	db, err := postgress.ConnectPostgres()
	if err != nil {
		log.Println("Error on connection", err)
		return nil
	}
	url := os.Getenv("kafka_url")
	log.Println("URl>>>>>>", url)
	// kafka, err  := kafka.ConnectKafka(os.Getenv("kafka_url"))
	// if err != nil {
	// 	log.Println("Error on kafka", err)
	// 	return nil
	// }

	rdb, err := rediss.ConnectRedis()
	if err != nil {
		log.Println("Error on connection redis...", err)
		return nil
	}
	return &UserServer{Db: db, Rdb: rdb}
}

func (s *UserServer) RegisterUserr(ctx context.Context, req *genuser.RegisterReq) (*genuser.RegisterResp, error) {
	log.Println("Keldi 1>>>>>>>>>>>>")
	user := models.RegisterReq{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Age:      req.Age,
	}
	redisKey := "user:" + user.Email
	userData := map[string]interface{}{
		"Name":     user.Name,
		"Email":    user.Email,
		"Password": user.Password,
		"Age":      user.Age,
	}

	err := s.Rdb.HSet(ctx, redisKey, userData).Err()
	log.Println("xatolik redisadami 1>>>>>>>>>>>>", err)
	if err != nil {
		log.Println("Failed to save user data to Redis:", err)
		return nil, err
	}
	


	err = s.Rdb.Set(ctx, "email:"+req.Email, req.Email, 150*time.Second).Err()
	log.Println("xatolik set redisadami 3>>>>>>>>>>>>", err)
	

	if err != nil {
		log.Println("Failed to save email to Redis:", err)
		return nil, err
	}

	Code, err := email.SendEmail(req.Email)
	log.Println("xatolik emaildami 4>>>>>>>>>>>>", err)
	if err != nil {
		return nil, errors.New("error on email write correct email")
	}

	err = s.Rdb.Set(ctx, "verification:"+req.Email, Code, 150*time.Second).Err()
	log.Println("xatolik verifixatisiondami 5>>>>>>>>>>>>", err)
	if err != nil {
		log.Println("Failed to save verification code to Redis:", err)
		return nil, err
	}
	log.Println("Xatolik kemadikkuuuu", err)
	return &genuser.RegisterResp{
		Message: "Check Your Email Please",
	}, nil
}

func (s *UserServer) VerifyUserr(ctx context.Context, req *genuser.VerifyReq) (*genuser.VerifyResp, error) {

	storedEmail, err := s.Rdb.Get(ctx, "email:"+req.Email).Result()
	log.Println("Stored email>>>>", storedEmail)
	if err != nil || storedEmail != req.Email {
		log.Println("smth wrong or match the email from Redis:", err)
		return &genuser.VerifyResp{
			Message: "Email does not match.",
		}, err
	}

	password, err := s.Rdb.HGet(ctx, "user:"+req.Email, "Password").Result()
	if err != nil {
		log.Println("wrpmgg to retrieve password from Redis:", err)
		return &genuser.VerifyResp{
			Message: "wrongg to retrieve password.",
		}, err
	}
	log.Println("password from Redis:", password)

	storedCode, err := s.Rdb.Get(ctx, "verification:"+req.Email).Result()
	log.Println("Stored code >>>>", storedCode)

	if err != nil {
		log.Println("Did not get code from Redis:", err)
		return &genuser.VerifyResp{
			Message: "Failed to verify user.",
		}, err
	}

	if storedCode != req.Code {
		log.Println("Verification code does not match for email:", req.Email)
		return &genuser.VerifyResp{
			Message: "Wrong verification code.",
		}, status.Error(codes.InvalidArgument, "Invalid code")
	}

	userData, err := s.Rdb.HGetAll(ctx, "user:"+req.Email).Result()
	if err != nil {
		log.Println("Failed to retrieve user data from Redis:", err)
		return &genuser.VerifyResp{
			Message: "Failed to verify user.",
		}, status.Error(codes.InvalidArgument, "failed to get verify user")
	}
	var existEmail string
	err = s.Db.QueryRow("SELECT email FROM users WHERE email = $1", req.Email).Scan(&existEmail)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error checking if email exists in PostgreSQL:", err)
		return &genuser.VerifyResp{
			Message: "Failed to verify user.",
		}, err
	}
	if existEmail != "" {

		return &genuser.VerifyResp{
			Message: "Email already exists. Please Register with another email.",
		}, err
	}

	passwordd, err := hash.HashPassword(password)
	if err != nil {
		log.Println("DId not hash password")
		return nil, err
	}

	var userID int32
	err = s.Db.QueryRow("INSERT INTO users (name, email, password, age) VALUES ($1, $2, $3, $4) RETURNING user_id",
		userData["Name"], userData["Email"], passwordd, userData["Age"]).Scan(&userID)
	if err != nil {
		log.Println("wronggg to save user to PostgreSQL:", err)
		return &genuser.VerifyResp{
			Message: "wrongg to save user.",
		}, err
	}

	token, err := jwt.CreateToken(req.Email)
	if err != nil {
		log.Println("Error on createtoken", err)
		return nil, err
	}
	// s.Kafka.ProduceRegistrationEmail(&genuser.NewUser{
	// 	UserId: userID,
	// 	Email: userData["Email"],
	// 	Name: userData["Name"],
	// })
	return &genuser.VerifyResp{
		UserId:  userID,
		Token:   token,
		Message: "User verified successfully.",
	}, nil
}

func (s *UserServer) Loginn(ctx context.Context, req *genuser.LoginReq) (*genuser.LoginResp, error) {
	query := `SELECT user_id, password FROM users WHERE email=$1`

	var user genuser.LoginResp
	var hashedPassword string

	err := s.Db.QueryRow(query, req.Email).Scan(&user.UserId, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("failed to login: %v", err)
	}

	if !hash.ComparePassword(hashedPassword, req.Password) {
		return nil, fmt.Errorf("invalid email or password")
	}

	token, err := jwt.CreateToken(req.Email)
	if err != nil {
		log.Println("Error creating JWT token:", err)
		return nil, err
	}

	return &genuser.LoginResp{
		UserId:  user.UserId,
		Token:   token,
		Message: "Successfully logged in",
	}, nil
}

func (s *UserServer) DeleteUserr(ctx context.Context, req *genuser.DeleteUserReq) (*genuser.DeleteUserResp, error) {

	var userExists bool
	queryCheck := `SELECT EXISTS(SELECT 1 FROM users WHERE user_id=$1)`
	err := s.Db.QueryRow(queryCheck, req.UserId).Scan(&userExists)
	if err != nil {
		log.Println("Error checking if user exists:", err)
		return &genuser.DeleteUserResp{
			Message: "Failed to check user existence.",
		}, err
	}

	if !userExists {
		return &genuser.DeleteUserResp{
			Message: "User not found.",
		}, nil
	}

	queryDelete := `DELETE FROM users WHERE user_id=$1`
	_, err = s.Db.Exec(queryDelete, req.UserId)
	if err != nil {
		log.Println("Error deleting user:", err)
		return &genuser.DeleteUserResp{
			Message: "Failed to delete user.",
		}, err
	}

	return &genuser.DeleteUserResp{
		Message: "User deleted successfully.",
	}, nil
}

func (s *UserServer) GetUserr(ctx context.Context, req *genuser.GetUserReq) (*genuser.GetUserResp, error) {
	var user models.GetUserResp

	query := `SELECT user_id, name, email, age FROM users WHERE user_id=$1`
	var name, email string
	var age int32

	err := s.Db.QueryRowContext(ctx, query, req.UserId).Scan(&user.UserID, &name, &email, &age)
	log.Println("Scan qivoganda nima chiqdi>>>", user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User with ID %d not found", req.UserId)
		}
		return nil, status.Errorf(codes.Internal, "Failed to retrieve user: %v", err)
	}

	registerResp := &genuser.RegisterReq{
		Name:  name,
		Email: email,
		Age:   age,
	}

	return &genuser.GetUserResp{
		UserId:   user.UserID,
		Response: []*genuser.RegisterReq{registerResp},
	}, nil
}
