package main

import (
  "context"
  "fmt"
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  sharedConst "github.com/arcorium/nexa/shared/constant"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "google.golang.org/grpc/metadata"
  "log"
  "os"
)

func main() {
  option := grpc.WithTransportCredentials(insecure.NewCredentials())
  conn, err := grpc.NewClient("localhost:8080", option)
  if err != nil {
    log.Printf("failed to connect to grpc server: %v", err)
    log.Printf("Trying to connect again")
  }

  token := os.Getenv("TOKEN")
  md := metadata.New(map[string]string{
    "authorization": fmt.Sprintf("%s %s", sharedConst.DEFAULT_ACCESS_TOKEN_SCHEME, token),
  })

  client := authZv1.NewRoleServiceClient(conn)
  mdCtx := metadata.NewOutgoingContext(context.Background(), md)
  create, err := client.GetUsers(mdCtx, &authZv1.GetUserRolesRequest{
    UserId:            types.MustCreateId().String(),
    IncludePermission: false,
  })
  if err != nil {
    log.Fatalln(err)
  }

  sharedUtil.DoNothing(create)
}
