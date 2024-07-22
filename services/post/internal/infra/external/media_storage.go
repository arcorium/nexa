package external

import (
  "context"
  storagev1 "github.com/arcorium/nexa/proto/gen/go/file_storage/v1"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/post/config"
  "nexa/services/post/internal/domain/external"
  "nexa/services/post/util"
)

func NewMediaStorage(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.IMediaStore {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-media-storage",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &mediaStorageClient{
    client: storagev1.NewFileStorageServiceClient(conn),
    tracer: util.GetTracer(),
    cb:     breaker,
  }
}

type mediaStorageClient struct {
  client storagev1.FileStorageServiceClient
  tracer trace.Tracer

  cb *gobreaker.CircuitBreaker
}

func (m *mediaStorageClient) GetUrls(ctx context.Context, fileIds ...types.Id) ([]string, error) {
  ctx, span := m.tracer.Start(ctx, "MediaStorage.GetUrls")
  defer span.End()

  req := storagev1.GetFileUrlRequest{
    FileIds: sharedUtil.CastSlice(fileIds, sharedUtil.ToString[types.Id]),
  }

  res, err := m.cb.Execute(func() (interface{}, error) {
    return m.client.GetUrls(ctx, &req)
  })

  if err != nil {
    err = util.CastBreakerError(err)
    spanUtil.RecordError(err, span)
    return nil, err
  }

  return res.([]string), nil
}
