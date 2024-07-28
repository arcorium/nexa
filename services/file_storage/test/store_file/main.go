package main

import (
  "context"
  storagev1 "github.com/arcorium/nexa/proto/gen/go/file_storage/v1"
  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"
  "log"
  "os"
)

func main() {
  data, err := os.ReadFile("test/25231.png")
  if err != nil {
    log.Fatalln(err)
  }

  conn, err := grpc.NewClient("localhost:8080", grpc.WithInsecure())
  if err != nil {
    log.Fatalln(err)
  }

  token := "bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJuZXhhLWF1dGgiLCJleHAiOjE3MjE5MjcyMTIsIm5iZiI6MTcyMTkyMzYxMiwiaWF0IjoxNzIxOTIzNjEyLCJqdGkiOiIzNWE5ODM3OS0zZDFjLTQyYzMtODFjMy03OTIwMjkwZWQxN2YiLCJjaWQiOiI3OWFhNWQ5MS1kNmI3LTRkM2YtOTc4OS1kYzc0ODY1MGVhMjEiLCJ1aWQiOiI5YTAyYzU3Ni1iNjQ1LTQxYWUtOWQ5My0zZTlmMGY0ZGE2MDciLCJ1c2VybmFtZSI6ImFyY29yaXVtIiwicm9sZXMiOlt7ImlkIjoiOWY3NTI0NmUtNThjZi00OGMxLTg4NmYtMTQ3MWRjNWU1ZTRjIiwicm9sZSI6IlVzZXIiLCJwZXJtcyI6WyJhdXRoejpnZXQ6cm9sZSIsImF1dGh6OmRlbGV0ZTp1c2VyOnJvbGUiLCJtYWlsOmdldDp0YWciLCJhdXRobjpnZXQ6dXNlciIsImF1dGhuOmdldDpjcmVkIiwiYXV0aG46dXBkYXRlOnVzZXIiLCJhdXRobjpkZWxldGU6dXNlciIsImF1dGhuOmxvZ291dDp1c2VyIiwiZmlsZTpnZXQ6ZmlsZSIsImZpbGU6Z2V0Om1kIiwiZmlsZTpzdG9yZTpmaWxlIiwiZmlsZTp1cGRhdGU6ZmlsZSIsImZpbGU6ZGVsZXRlOmZpbGUiXX1dfQ.V46x_TqFJf9XZ51mpA1KU-TXQlnQL_sMc1QcjdkewRimvQ5mK0OhN35PnjgdQiYoVjAuxi1tHuTEp9xdLu59-u4VAdaJnas7lSzJK8QHWf0Y0vOAOrYmvMaYMuaPvSf39HSRZWWJlGPqbkXMht-SeWz5GDEsfa7MKXXHsF5OC_KL36CB_V6xrFZwT-UrSm1E5e4xbZzsmVBtAXwW-GXgMslohvz37tQ9-3JDYZWnDOnyEFFzUFoA7h9IUPrHOJQaMKQbgRtrof8PAv_HZjmbsF5o8h87pOM1llA442F7l5IqMaBcCtnwCaIzstV0ZJkkP1737ua1XvfAG-AJdygBmg"

  md := metadata.New(map[string]string{
    "authorization": token,
    "Content-Type":  "application/grpc",
  })
  ctx := metadata.NewOutgoingContext(context.Background(), md)

  client := storagev1.NewFileStorageServiceClient(conn)
  store, err := client.Store(ctx)
  if err != nil {
    log.Fatalln(err)
  }
  err = store.Send(&storagev1.StoreFileRequest{
    Filename: "25231.png",
    IsPublic: true,
    Chunk:    data,
  })
  if err != nil {
    log.Fatalln(err)
  }

  recv, err := store.CloseAndRecv()
  if err != nil {
    log.Fatalln(err)
  }
  log.Println(recv)
}
