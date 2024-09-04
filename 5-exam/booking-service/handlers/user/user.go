package userr

import (
	"booking-service/proto/genuser"
	"context"
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func HandleUser(id int) (int, error) {
	conn, err := grpc.NewClient("user_service:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return 0, err
	}
	client := genuser.NewUserServiceClient(conn)
	user, err := client.GetUser(context.Background(), &genuser.GetUserReq{UserId: int32(id)})
	if err != nil {
		log.Println("Error on connection handuser>>>", err)
		return 0, err
	}


	return int(user.UserId), nil
}
