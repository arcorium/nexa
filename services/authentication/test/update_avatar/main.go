package main

import (
  "context"
  authNv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"
  "io"
  "log"
  "os"
)

func main() {
  conn, err := grpc.NewClient("localhost:8080", grpc.WithInsecure())
  if err != nil {
    log.Fatalln(err)
  }

  token := "bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJuZXhhLWF1dGgiLCJleHAiOjE3MjIwMTIwNzksIm5iZiI6MTcyMTkyNTY3OSwiaWF0IjoxNzIxOTI1Njc5LCJqdGkiOiI1NjkxNTQ4ZS1mNWU0LTRlZjQtYjYxMS0yMmE2YTRmODE1NTkiLCJjaWQiOiIzNzQ3NjgxOC1iMTBkLTQyNGMtYTU1ZS0zZDQ2MzNkNjNiNWUiLCJ1aWQiOiI2ZDRkZDQwNi0wMTc5LTRmZjAtYTdlMi00OWEwYmM2Nzg4ZTkiLCJ1c2VybmFtZSI6ImFyY29yaXVtIiwicm9sZXMiOlt7ImlkIjoiNTdlNzQ1NmYtNjIxMC00NjEyLTk3MDYtZDIyNGExYzE5Y2FkIiwicm9sZSI6IlVzZXIiLCJwZXJtcyI6WyJhdXRoejpnZXQ6cm9sZSIsImF1dGh6OmRlbGV0ZTp1c2VyOnJvbGUiLCJmaWxlOmdldDpmaWxlIiwiZmlsZTpnZXQ6bWQiLCJmaWxlOnN0b3JlOmZpbGUiLCJmaWxlOnVwZGF0ZTpmaWxlIiwiZmlsZTpkZWxldGU6ZmlsZSIsIm1haWw6Z2V0OnRhZyIsImF1dGhuOmdldDp1c2VyIiwiYXV0aG46Z2V0OmNyZWQiLCJhdXRobjp1cGRhdGU6dXNlciIsImF1dGhuOmRlbGV0ZTp1c2VyIiwiYXV0aG46bG9nb3V0OnVzZXIiXX1dfQ.V1zuPaqe3IqrIxFStr1wWHeaWGEtdMWfitmHv5wGdnFZIxGrUV9JF00PzHw4Sgx059ohVuZSzPUj8mOna4nDvVn_uO9nnTsAKbYg_MdC8sF995ZmsUDt26cIbwF7ip02VOSfRBKk91N2D3B0ar_e7o0hUwRcMebOHMCvmZi8c-u-hKLOjIE0Se4-Fe-tnDv7RNMXP9vxdHiEt_5mno8m--QKYrlyIgE4uaqFmJQ8kGn_AgC6no-YVSVc2YVO8aY_etLAP7gilfHdP_3gqzoe8KsUwWhRwsIUGsHiA7kQd9FyOID7V2b07YltR_QioxYszamBPRn3psDtxiMb944UMg"

  md := metadata.New(map[string]string{
    "authorization": token,
    "Content-Type":  "application/grpc",
  })
  ctx := metadata.NewOutgoingContext(context.Background(), md)

  client := authNv1.NewUserServiceClient(conn)
  sv, err := client.UpdateAvatar(ctx)
  if err != nil {
    log.Fatalln(err)
  }

  filename := "test/25231.png"
  data, err := os.ReadFile(filename)
  if err != nil {
    log.Fatalln(err)
  }

  for {
    err = sv.Send(&authNv1.UpdateProfileAvatarRequest{
      Filename: "25231.png",
      Chunk:    data,
    })
    if err != nil {
      if err != io.EOF {
        log.Fatalln(err)
      }
      sv.CloseSend()
      log.Fatalln(err)
    }
    break
  }

  _, err = sv.CloseAndRecv()
  if err != nil {
    log.Fatalln(err)
  }
}
