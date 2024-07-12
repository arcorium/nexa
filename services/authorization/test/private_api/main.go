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
)

func main() {
  option := grpc.WithTransportCredentials(insecure.NewCredentials())
  conn, err := grpc.NewClient("localhost:8080", option)
  if err != nil {
    log.Printf("failed to connect to grpc server: %v", err)
    log.Printf("Trying to connect again")
  }

  token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJuZXhhLXRva2VuX2dlbmVyYXRvciIsInN1YiI6InNldHVwIiwiZXhwIjoxNzIwODMwNjc1LCJpYXQiOjE3MjA3OTQ2NzUsImp0aSI6IjlleHlnNW05cGZIV3E4N0hPSHdqNWxMYVpNeWRtVk5BIn0.l5YKW-xnUDzcam9hRe0-nw_hNLjBgqoA7-tbnHbOfRT5rEhqzs_9yGNKlp0ugL_UZNLuZjfB4VFGKOy1uxx6yWNxXHHw_nKT9mOtEgdmaeX4DJ20sSo3jPTaMWp9zEPjcHZ7jRsePe19itWMY-Ugg1SwW6_AQVk0k37u-nFFczbT88OJp0V4RMmLC_FKurRP8imb8C_E1OI0FtUDeGZLk1HM2Kz8miTH-NZqjqQZZwaT4qcD67946wSSSIkZaihmRr1AB6jSpZAPZsby_ccghV4hAWQlvAmKHau7SKYS0hU0LrxrW3rm8W8w9QwUq2fpZZvf5OKseP43oYz6j-t-vQ"
  md := metadata.New(map[string]string{
    "authorization": fmt.Sprintf("%s %s", sharedConst.DEFAULT_ACCESS_TOKEN_SCHEME, token),
  })

  client := authZv1.NewPermissionServiceClient(conn)
  mdCtx := metadata.NewOutgoingContext(context.Background(), md)
  create, err := client.Seed(mdCtx, &authZv1.SeedPermissionRequest{
    Permissions: []*authZv1.CreatePermissionRequest{
      {
        Resource: "origin",
        Action:   "uwa",
      },
    },
  })
  if err != nil {
    log.Fatalln(err)
  }

  sharedUtil.DoNothing(create)
}
