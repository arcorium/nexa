package main

import (
  "context"
  "fmt"
  authZv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  sharedConst "github.com/arcorium/nexa/shared/constant"
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

  client := authZv1.NewPermissionServiceClient(conn)
  mdCtx := metadata.NewOutgoingContext(context.Background(), md)
  create, err := client.Create(mdCtx, &authZv1.CreatePermissionRequest{
    Resource: "something",
    Action:   "doin",
  })
  if err != nil {
    log.Fatalln(err)
  }

  sharedUtil.DoNothing(create)
}
