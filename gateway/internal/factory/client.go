package factory

import (
  "context"
  "errors"
  "fmt"
  authenticationv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  authorizationv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  commentv1 "github.com/arcorium/nexa/proto/gen/go/comment/v1"
  feedv1 "github.com/arcorium/nexa/proto/gen/go/feed/v1"
  storagev1 "github.com/arcorium/nexa/proto/gen/go/file_storage/v1"
  mailerv1 "github.com/arcorium/nexa/proto/gen/go/mailer/v1"
  postv1 "github.com/arcorium/nexa/proto/gen/go/post/v1"
  reactionv1 "github.com/arcorium/nexa/proto/gen/go/reaction/v1"
  relationv1 "github.com/arcorium/nexa/proto/gen/go/relation/v1"
  "github.com/arcorium/nexa/shared/logger"
  "google.golang.org/grpc"
  "google.golang.org/grpc/health/grpc_health_v1"
  "nexa/gateway/config"
)

type Client struct {
  Authn authenticationv1.AuthenticationServiceClient
  User  authenticationv1.UserServiceClient

  Authz      authorizationv1.AuthorizationServiceClient
  Permission authorizationv1.PermissionServiceClient
  Role       authorizationv1.RoleServiceClient

  Comment commentv1.CommentServiceClient
  Feed    feedv1.FeedServiceClient
  Storage storagev1.FileStorageServiceClient

  Mailer    mailerv1.MailerServiceClient
  MailerTag mailerv1.TagServiceClient

  Post     postv1.PostServiceClient
  Reaction reactionv1.ReactionServiceClient

  Follow relationv1.FollowServiceClient
  Block  relationv1.BlockServiceClient

  //Token tokenv1.TokenServiceClient

  conn Connection
}

type Connection struct {
  Authn    *grpc.ClientConn
  Authz    *grpc.ClientConn
  Comment  *grpc.ClientConn
  Feed     *grpc.ClientConn
  Storage  *grpc.ClientConn
  Mailer   *grpc.ClientConn
  Post     *grpc.ClientConn
  Reaction *grpc.ClientConn
  Relation *grpc.ClientConn
  //Token    *grpc.ClientConn
}

var ErrRetryFailed = errors.New("reached max retry, failed to do the actions")

func (c *Connection) retry(ctx context.Context, maxRetry uint, f func(context.Context) error) error {
  var err error
  for range maxRetry {
    err = f(ctx)
    if err == nil {
      return nil
    }
  }
  return ErrRetryFailed
}

func (c *Connection) checkHealths(ctx context.Context) error {
  conns := map[string]*grpc.ClientConn{
    "Authentication": c.Authn,
    "Authorization":  c.Authz,
    "Comment":        c.Comment,
    "Feed":           c.Feed,
    "Storage":        c.Storage,
    "Mailer":         c.Mailer,
    "Post":           c.Post,
    "Reaction":       c.Reaction,
    "Relation":       c.Relation,
    //"Token":          c.Token,
  }
  for name, conn := range conns {
    client := grpc_health_v1.NewHealthClient(conn)
    err := c.retry(ctx, 5, func(ctx context.Context) error {
      check, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
      if err != nil {
        return err
      }
      if check.Status != grpc_health_v1.HealthCheckResponse_SERVING {
        return fmt.Errorf("service %s failed to serve: %s", name, check.Status.String())
      }
      return nil
    })
    if err != nil {
      return fmt.Errorf("service %s health check failed: %s", name, err)
    }
    logger.Infof("Service %s health check pass", name)
  }
  return nil
}

func NewClient(conf *config.Service, opts ...grpc.DialOption) (*Client, error) {
  client := &Client{}

  {
    conn, err := grpc.NewClient(conf.Authn, opts...)
    if err != nil {
      return nil, err
    }
    client.Authn = authenticationv1.NewAuthenticationServiceClient(conn)
    client.User = authenticationv1.NewUserServiceClient(conn)
    client.conn.Authn = conn
  }
  {
    conn, err := grpc.NewClient(conf.Authz, opts...)
    if err != nil {
      return nil, err
    }
    client.Authz = authorizationv1.NewAuthorizationServiceClient(conn)
    client.Permission = authorizationv1.NewPermissionServiceClient(conn)
    client.Role = authorizationv1.NewRoleServiceClient(conn)

    client.conn.Authz = conn
  }
  {
    conn, err := grpc.NewClient(conf.Comment, opts...)
    if err != nil {
      return nil, err
    }
    client.Comment = commentv1.NewCommentServiceClient(conn)

    client.conn.Comment = conn
  }
  {
    conn, err := grpc.NewClient(conf.Feed, opts...)
    if err != nil {
      return nil, err
    }
    client.Feed = feedv1.NewFeedServiceClient(conn)

    client.conn.Feed = conn
  }
  {
    conn, err := grpc.NewClient(conf.Storage, opts...)
    if err != nil {
      return nil, err
    }
    client.Storage = storagev1.NewFileStorageServiceClient(conn)

    client.conn.Storage = conn
  }
  {
    conn, err := grpc.NewClient(conf.Mailer, opts...)
    if err != nil {
      return nil, err
    }
    client.Mailer = mailerv1.NewMailerServiceClient(conn)
    client.MailerTag = mailerv1.NewTagServiceClient(conn)

    client.conn.Mailer = conn
  }
  {
    conn, err := grpc.NewClient(conf.Post, opts...)
    if err != nil {
      return nil, err
    }
    client.Post = postv1.NewPostServiceClient(conn)

    client.conn.Post = conn
  }
  {
    conn, err := grpc.NewClient(conf.Reaction, opts...)
    if err != nil {
      return nil, err
    }
    client.Reaction = reactionv1.NewReactionServiceClient(conn)

    client.conn.Reaction = conn
  }
  {
    conn, err := grpc.NewClient(conf.Relation, opts...)
    if err != nil {
      return nil, err
    }
    client.Follow = relationv1.NewFollowServiceClient(conn)
    client.Block = relationv1.NewBlockServiceClient(conn)

    client.conn.Relation = conn
  }
  //{
  //  conn, err := grpc.NewClient(conf.Token, opts...)
  //  if err != nil {
  //    return nil, err
  //  }
  //  client.Token = tokenv1.NewTokenServiceClient(conn)
  //
  //  client.conn.Token = conn
  //}

  return client, client.conn.checkHealths(context.Background())
}
