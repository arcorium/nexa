package factory

import (
  "github.com/arcorium/nexa/shared/logger"
  "google.golang.org/grpc"
  "nexa/services/authentication/config"
  "nexa/services/authentication/internal/domain/external"
  . "nexa/services/authentication/internal/infra/external"
)

type External struct {
  Token   external.ITokenClient
  Mail    external.IMailerClient
  Storage external.IFileStorageClient
  Role    external.IRoleClient

  connections []*grpc.ClientConn
}

func (e *External) Close() {
  for _, conn := range e.connections {
    if err := conn.Close(); err != nil {
      logger.Warnf("Error on closing external: %s", err.Error())
    }
  }
}

func NewExternalWithConn(tokenConn, mailConn, storageConn, roleConn *grpc.ClientConn, conf *config.CircuitBreaker) *External {
  return &External{
    Token:       NewTokenClient(tokenConn, conf),
    Mail:        NewMailerClient(mailConn, conf),
    Storage:     NewFileStorageClient(storageConn, conf),
    Role:        NewRoleClient(roleConn, conf),
    connections: []*grpc.ClientConn{tokenConn, mailConn, storageConn, roleConn},
  }
}

func NewExternalWithConfig(conf *config.Service, breaker *config.CircuitBreaker, options ...grpc.DialOption) (*External, error) {
  roleConn, err := grpc.NewClient(conf.Authorization, options...)
  if err != nil {
    return nil, err
  }

  tokenConn, err := grpc.NewClient(conf.Token, options...)
  if err != nil {
    return nil, err
  }

  mailConn, err := grpc.NewClient(conf.Mailer, options...)
  if err != nil {
    return nil, err
  }

  storageConn, err := grpc.NewClient(conf.FileStorage, options...)
  if err != nil {
    return nil, err
  }
  return NewExternalWithConn(tokenConn, mailConn, storageConn, roleConn, breaker), nil
}
