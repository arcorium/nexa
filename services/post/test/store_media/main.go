package main

import (
  "context"
  storagev1 "github.com/arcorium/nexa/proto/gen/go/file_storage/v1"
  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"
  "log"
  "os"
  "path"
)

func sendFile(filepath string, md metadata.MD, client storagev1.FileStorageServiceClient) error {
  _, filename := path.Split(filepath)
  data, err := os.ReadFile(filepath)
  if err != nil {
    return err
  }

  ctx := metadata.NewOutgoingContext(context.Background(), md)

  store, err := client.Store(ctx)
  if err != nil {
    return err
  }
  err = store.Send(&storagev1.StoreFileRequest{
    Filename: filename,
    IsPublic: true,
    Chunk:    data,
  })
  if err != nil {
    return err
  }

  recv, err := store.CloseAndRecv()
  if err != nil {
    return err
  }
  log.Println(recv)
  return nil
}

func main() {
  conn, err := grpc.NewClient("localhost:8083", grpc.WithInsecure())
  if err != nil {
    log.Fatalln(err)
  }

  token := "bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJuZXhhLWF1dGgiLCJleHAiOjE3MjIxMDgzNjksIm5iZiI6MTcyMjAyMTk2OSwiaWF0IjoxNzIyMDIxOTY5LCJqdGkiOiJjMDQ0OGJlNy03NGUwLTQ0M2QtYmM3Mi0yOTExMGExZjkxNWYiLCJjaWQiOiJkYWYwZmVmZi1mZTRkLTQ1ZDQtOTBjZi1hN2U2NTU5MWI5N2UiLCJ1aWQiOiJlNjhkNzQ4My1jMGY1LTRhY2ItOGMxNC05OTMxYTQ2MGQwZGEiLCJ1c2VybmFtZSI6ImFyY29yaXVtIiwicm9sZXMiOlt7ImlkIjoiOTFiYTgzYTAtMDg4OC00NWQxLTg5ZDUtNzAzN2E0ZTg2YWFlIiwicm9sZSI6IlVzZXIiLCJwZXJtcyI6WyJhdXRoejpnZXQ6cm9sZSIsImF1dGh6OmRlbGV0ZTp1c2VyOnJvbGUiLCJwb3N0OmNyZWF0ZTpwb3N0IiwicG9zdDp1cGRhdGU6cG9zdCIsInBvc3Q6Z2V0OnBvc3QiLCJwb3N0OmRlbGV0ZTpwb3N0IiwiZmlsZTpnZXQ6ZmlsZSIsImZpbGU6Z2V0Om1kIiwiZmlsZTpzdG9yZTpmaWxlIiwiZmlsZTp1cGRhdGU6ZmlsZSIsImZpbGU6ZGVsZXRlOmZpbGUiLCJhdXRobjpnZXQ6dXNlciIsImF1dGhuOmdldDpjcmVkIiwiYXV0aG46dXBkYXRlOnVzZXIiLCJhdXRobjpkZWxldGU6dXNlciIsImF1dGhuOmxvZ291dDp1c2VyIiwiYXV0aG46cmVxOnVzZXI6dmVyaWYiLCJtYWlsOmdldDp0YWciLCJyZWxhdGlvbjpnZXQ6Zm9sbG93IiwicmVsYXRpb246ZGVsZXRlOmZvbGxvdyIsInJlbGF0aW9uOmNyZWF0ZTpibG9jayIsInJlbGF0aW9uOmdldDpibG9jayIsInJlbGF0aW9uOmNyZWF0ZTpmb2xsb3ciLCJyZWxhdGlvbjpkZWxldGU6YmxvY2siXX1dfQ.c9YJw8RJJC-BUtltdq39lySzkWaj77y1lAXS5ybDduPvHZctGJJe6Oc9HTKvjsh2t2Kk-c-nB124nK8ZzxO4T3dJpgqvemRqbGUpR8cWf1gOGBhc25Nc5p1n5lcK0_w5DY4sOWpFZDIarufSXo-Rb2PoDmb8qnNZchArHRs-GiUmvnXu9EAuK2oINbRTVWiX9qhvt4KbutiAyfFJFSjqSocEkHo7A8hDNqUpQ8hoj6nlo7X3-gbEpoBan0detX4lLQt-L0xFkkys-XIU_tmZ9saBN9CYuGCNai2WZwd6eau63IUl85liDLt2O6wP0ZuMOBdUq5vmg6PpUIulTQBtCQ"

  md := metadata.New(map[string]string{
    "authorization": token,
    "Content-Type":  "application/grpc",
  })

  client := storagev1.NewFileStorageServiceClient(conn)

  err = sendFile("test/240.png", md, client)
  if err != nil {
    log.Fatalln(err)
  }

  err = sendFile("test/25231.png", md, client)
  if err != nil {
    log.Fatalln(err)
  }

  err = sendFile("test/ForBiggerMeltdowns.mp4", md, client)
  if err != nil {
    log.Fatalln(err)
  }
}
