package factory

import (
  "github.com/arcorium/nexa/shared/logger"
  "google.golang.org/grpc"
  "nexa/services/post/config"
  "nexa/services/post/internal/domain/external"
  . "nexa/services/post/internal/infra/external"
)

type External struct {
  Comment      external.ICommentClient
  Follow       external.IFollowClient
  MediaStorage external.IMediaStore
  Reaction     external.IReactionClient
  User         external.IUserClient

  connections []*grpc.ClientConn
}

func (e *External) Close() {
  for _, conn := range e.connections {
    if err := conn.Close(); err != nil {
      logger.Warnf("Error on closing external: %s", err.Error())
    }
  }
}

func NewExternalWithConn(commentConn, followConn, mediaConn, reactionConn, userConn *grpc.ClientConn, breaker *config.CircuitBreaker) *External {
  return &External{
    Comment:      NewComment(commentConn, breaker),
    Follow:       NewFollow(followConn, breaker),
    MediaStorage: NewMediaStorage(mediaConn, breaker),
    Reaction:     NewReaction(reactionConn, breaker),
    User:         NewUser(userConn, breaker),
    connections:  []*grpc.ClientConn{commentConn, followConn, mediaConn, reactionConn, userConn},
  }
}

func NewExternalWithConfig(conf *config.Service, breaker *config.CircuitBreaker, options ...grpc.DialOption) (*External, error) {
  commentConn, err := grpc.NewClient(conf.Comment, options...)
  if err != nil {
    return nil, err
  }

  followConn, err := grpc.NewClient(conf.Follow, options...)
  if err != nil {
    return nil, err
  }

  mediaConn, err := grpc.NewClient(conf.MediaStorage, options...)
  if err != nil {
    return nil, err
  }

  reactionConn, err := grpc.NewClient(conf.Reaction, options...)
  if err != nil {
    return nil, err
  }

  userConn, err := grpc.NewClient(conf.User, options...)
  if err != nil {
    return nil, err
  }

  return NewExternalWithConn(commentConn, followConn, mediaConn, reactionConn, userConn, breaker), nil
}
