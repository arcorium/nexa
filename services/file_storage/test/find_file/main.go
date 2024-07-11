package main

import (
  "context"
  storagev1 "github.com/arcorium/nexa/proto/gen/go/file_storage/v1"
  "google.golang.org/grpc"
  "io"
  "log"
  "os"
)

func main() {
  conn, err := grpc.NewClient("localhost:8080", grpc.WithInsecure())
  if err != nil {
    log.Fatalln(err)
  }

  client := storagev1.NewFileStorageServiceClient(conn)
  find, err := client.Find(context.Background(), &storagev1.FindFileRequest{FileId: "93fa62d3-a0de-4b98-8a47-0430d327b7f8"})
  if err != nil {
    log.Fatalln(err)
  }

  bytes := make([]byte, 0)
  var filename string

  for {
    recv, err := find.Recv()
    if err != nil {
      if err != io.EOF {
        log.Fatalln(err)
      }
      break
    }
    bytes = append(bytes, recv.Chunk...)
    filename = recv.Filename
  }

  err = os.WriteFile("test/got_"+filename, bytes, 777)
  if err != nil {
    log.Fatalln(err)
  }
}
