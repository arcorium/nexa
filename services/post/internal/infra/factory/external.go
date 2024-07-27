package factory

import (
  "github.com/arcorium/nexa/shared/logger"
  "google.golang.org/grpc"
  "nexa/services/post/config"
  "nexa/services/post/internal/domain/external"
  . "nexa/services/post/internal/infra/external"
)

type External struct {
  Relation     external.IRelationClient
  MediaStorage external.IMediaStoreClient

  connections []*grpc.ClientConn
}

func (e *External) Close() {
  for _, conn := range e.connections {
    if err := conn.Close(); err != nil {
      logger.Warnf("Error on closing external: %s", err.Error())
    }
  }
}

func NewExternalWithConn(relationConn, mediaConn *grpc.ClientConn, breaker *config.CircuitBreaker) *External {
  return &External{
    Relation:     NewRelationClient(relationConn, breaker),
    MediaStorage: NewMediaStorage(mediaConn, breaker),
    connections:  []*grpc.ClientConn{relationConn, mediaConn},
  }
}

func NewExternalWithConfig(conf *config.Service, breakerConf *config.CircuitBreaker, options ...grpc.DialOption) (*External, error) {
  relationConn, err := grpc.NewClient(conf.Relation, options...)
  if err != nil {
    return nil, err
  }

  mediaConn, err := grpc.NewClient(conf.MediaStorage, options...)
  if err != nil {
    return nil, err
  }

  return NewExternalWithConn(relationConn, mediaConn, breakerConf), nil
}
