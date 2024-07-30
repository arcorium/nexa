package factory

import (
  "github.com/arcorium/nexa/shared/logger"
  "github.com/redis/go-redis/v9"
  "google.golang.org/grpc"
  "nexa/services/feed/config"
  "nexa/services/feed/internal/domain/external"
  . "nexa/services/feed/internal/infra/external"
)

type External struct {
  Relation external.IRelationClient
  Reaction external.IReactionClient
  Comment  external.ICommentClient
  Post     external.IPostClient

  connections []*grpc.ClientConn
}

func (e *External) Close() {
  for _, conn := range e.connections {
    if err := conn.Close(); err != nil {
      logger.Warnf("Error on closing external: %s", err.Error())
    }
  }
}

func NewExternalWithConn(relationConn, reactionConn, commentConn, postConn *grpc.ClientConn, redis redis.UniversalClient, conf *config.CircuitBreaker) *External {
  ext := &External{
    Relation:    NewRelationClient(relationConn, redis, conf),
    Reaction:    NewReactionClient(reactionConn, conf),
    Comment:     NewCommentClient(commentConn, conf),
    connections: []*grpc.ClientConn{relationConn, reactionConn, commentConn, postConn},
  }
  ext.Post = NewPostClient(postConn, redis, ext.Reaction, ext.Comment, conf)

  return ext
}

func NewExternalWithConfig(conf *config.Service, redis redis.UniversalClient, breaker *config.CircuitBreaker, options ...grpc.DialOption) (*External, error) {
  relConn, err := grpc.NewClient(conf.Relation, options...)
  if err != nil {
    return nil, err
  }

  reactConn, err := grpc.NewClient(conf.Reaction, options...)
  if err != nil {
    return nil, err
  }

  commentConn, err := grpc.NewClient(conf.Comment, options...)
  if err != nil {
    return nil, err
  }

  postConn, err := grpc.NewClient(conf.Post, options...)
  if err != nil {
    return nil, err
  }
  return NewExternalWithConn(relConn, reactConn, commentConn, postConn, redis, breaker), nil
}
