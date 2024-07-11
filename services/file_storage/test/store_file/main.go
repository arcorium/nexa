package main

import (
  "context"
  storagev1 "github.com/arcorium/nexa/proto/gen/go/file_storage/v1"
  "google.golang.org/grpc"
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

  client := storagev1.NewFileStorageServiceClient(conn)
  store, err := client.Store(context.Background())
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
