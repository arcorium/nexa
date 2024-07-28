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

  token := "bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJuZXhhLWF1dGgiLCJleHAiOjE3MjIyNjU4OTQsIm5iZiI6MTcyMjE3OTQ5NCwiaWF0IjoxNzIyMTc5NDk0LCJqdGkiOiJhZjNjZThjNS04NzkwLTRhMWUtOWY5Zi1hMjE1OGYyY2I0N2YiLCJjaWQiOiJmZjBiZWJiMC04YTBhLTRiYjMtYWE3NC01NWYwZTZkNjQyZGEiLCJ1aWQiOiI0MjY2Mjc0Yi1jYTU3LTQ2YjgtYmIyMy0wMTZiMGM2MjBjNTciLCJ1c2VybmFtZSI6ImFyY29yaXVtIiwicm9sZXMiOlt7ImlkIjoiNTNjNDg4MzAtMGM2Zi00MGFmLThkNjQtZjQzODNkYTZmODc0Iiwicm9sZSI6IlVzZXIiLCJwZXJtcyI6WyJhdXRoejpnZXQ6cm9sZSIsImF1dGh6OmRlbGV0ZTp1c2VyOnJvbGUiLCJsaWtlOmNyZWF0ZTpyZWFjdGlvbiIsImxpa2U6Z2V0OnJlYWN0aW9uIiwibGlrZTpkZWxldGU6cmVhY3Rpb24iLCJjb21tZW50OmNyZWF0ZTpjb21tZW50IiwiY29tbWVudDp1cGRhdGU6Y29tbWVudCIsImNvbW1lbnQ6Z2V0OmNvbW1lbnQiLCJjb21tZW50OmRlbGV0ZTpjb21tZW50IiwicG9zdDpjcmVhdGU6cG9zdCIsInBvc3Q6dXBkYXRlOnBvc3QiLCJwb3N0OmdldDpwb3N0IiwicG9zdDpkZWxldGU6cG9zdCIsIm1haWw6Z2V0OnRhZyIsImZpbGU6Z2V0OmZpbGUiLCJmaWxlOmdldDptZCIsImZpbGU6c3RvcmU6ZmlsZSIsImZpbGU6dXBkYXRlOmZpbGUiLCJmaWxlOmRlbGV0ZTpmaWxlIiwicmVsYXRpb246Z2V0OmZvbGxvdyIsInJlbGF0aW9uOmRlbGV0ZTpmb2xsb3ciLCJyZWxhdGlvbjpjcmVhdGU6YmxvY2siLCJyZWxhdGlvbjpnZXQ6YmxvY2siLCJyZWxhdGlvbjpjcmVhdGU6Zm9sbG93IiwicmVsYXRpb246ZGVsZXRlOmJsb2NrIiwiYXV0aG46Z2V0OnVzZXIiLCJhdXRobjpnZXQ6Y3JlZCIsImF1dGhuOnVwZGF0ZTp1c2VyIiwiYXV0aG46ZGVsZXRlOnVzZXIiLCJhdXRobjpsb2dvdXQ6dXNlciIsImF1dGhuOnJlcTp1c2VyOnZlcmlmIl19XX0.UhDkmgiU64IHcyz2X-UMvBBEw5rPEE_CmFvNPUDsKE3Cid8Z6dNiJR2ftesVZ50J8AwqH9hc5okh9mqhz4noA_g1CZkKfSiqCbajD-gx4vuHQ5qAQGX_wz5aaG6kUDMhrohNIcep87fMJWDvTszdiuuiSk-BSg4DydcAKRQNLu-uZigBH55jZDraj2MXXGEYRBi3RQJ0ozZ-KAqe4T-ZC2Fbh2GRHLsih_4Bn0pDC5tkVFazxQ65SFO49lj5HX01Yv_qGvOJZmQkzvPLUIiQ0-3bX9xtxrQufOqiAmAWtWhRu9tD_Sr18t9W9B71HmofzZVuOr57s1PsnoxZfvYz4g"

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
