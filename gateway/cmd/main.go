package main

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	v1 "nexa/proto/generated/golang/user/v1"
	"nexa/shared/util"
	"os"
)

func main() {
	cred := insecure.NewCredentials()
	conn, err := grpc.NewClient("localhost:9999", grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatalln(err)
	}

	profileService := v1.NewProfileServiceClient(conn)

	//claims := &token.UserClaims{
	//  UserId: uuid.NewString(),
	//  Roles:  []string{"user"},
	//}
	//
	resp, err := profileService.UpdateAvatar(context.Background())

	bytes, _ := os.ReadFile("/mnt/data/Download/bukti_pemaparan_UAS.pdf")

	for {
		err = resp.Send(&v1.UpdateProfileAvatarRequest{
			UserId:   uuid.NewString(),
			Filename: "bukti_pemaparan_UAS.pdf",
			Chunk:    bytes,
		})

		break
	}

	stat, err := resp.CloseAndRecv()

	util.DoNothing(resp, stat)
}
