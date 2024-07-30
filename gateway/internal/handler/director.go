package handler

import (
  "context"
  authenticationv1 "github.com/arcorium/nexa/proto/gen/go/authentication/v1"
  authorizationv1 "github.com/arcorium/nexa/proto/gen/go/authorization/v1"
  commentv1 "github.com/arcorium/nexa/proto/gen/go/comment/v1"
  feedv1 "github.com/arcorium/nexa/proto/gen/go/feed/v1"
  storagev1 "github.com/arcorium/nexa/proto/gen/go/file_storage/v1"
  mailerv1 "github.com/arcorium/nexa/proto/gen/go/mailer/v1"
  postv1 "github.com/arcorium/nexa/proto/gen/go/post/v1"
  reactionv1 "github.com/arcorium/nexa/proto/gen/go/reaction/v1"
  relationv1 "github.com/arcorium/nexa/proto/gen/go/relation/v1"
  tokenv1 "github.com/arcorium/nexa/proto/gen/go/token/v1"
  "github.com/arcorium/nexa/shared/logger"
  "github.com/siderolabs/grpc-proxy/proxy"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  "nexa/gateway/config"
  "nexa/gateway/internal/factory"
  "strings"
)

func Director(backend *factory.Backend) proxy.StreamDirector {
  return func(ctx context.Context, fullMethodName string) (proxy.Mode, []proxy.Backend, error) {
    split := strings.Split(fullMethodName, "/")
    if len(split) > 3 {
      return proxy.One2One, nil, status.Errorf(codes.Unimplemented, "Unknown method")
    }
    serviceName := split[1]
    methodName := split[2]

    // Make sure we never forward internal services.
    if strings.HasPrefix(serviceName, "nexa.internal.") {
      return proxy.One2One, nil, status.Errorf(codes.Unimplemented, "Unknown method %s", methodName)
    }

    if config.IsDebug() {
      logger.Infof("Forwarding %s", fullMethodName)
    }

    switch serviceName {
    case authenticationv1.AuthenticationService_ServiceDesc.ServiceName:
      fallthrough
    case authenticationv1.UserService_ServiceDesc.ServiceName:
      return proxy.One2One, []proxy.Backend{backend.Authn()}, nil
    case authorizationv1.AuthorizationService_ServiceDesc.ServiceName:
      fallthrough
    case authorizationv1.PermissionService_ServiceDesc.ServiceName:
      fallthrough
    case authorizationv1.RoleService_ServiceDesc.ServiceName:
      return proxy.One2One, []proxy.Backend{backend.Authz()}, nil
    case commentv1.CommentService_ServiceDesc.ServiceName:
      return proxy.One2One, []proxy.Backend{backend.Comment()}, nil
    case feedv1.FeedService_ServiceDesc.ServiceName:
      return proxy.One2One, []proxy.Backend{backend.Feed()}, nil
    case storagev1.FileStorageService_ServiceDesc.ServiceName:
      return proxy.One2One, []proxy.Backend{backend.Storage()}, nil
    case mailerv1.MailerService_ServiceDesc.ServiceName:
      fallthrough
    case mailerv1.TagService_ServiceDesc.ServiceName:
      return proxy.One2One, []proxy.Backend{backend.Mailer()}, nil
    case postv1.PostService_ServiceDesc.ServiceName:
      return proxy.One2One, []proxy.Backend{backend.Post()}, nil
    case reactionv1.ReactionService_ServiceDesc.ServiceName:
      return proxy.One2One, []proxy.Backend{backend.Reaction()}, nil
    case relationv1.FollowService_ServiceDesc.ServiceName:
      fallthrough
    case relationv1.BlockService_ServiceDesc.ServiceName:
      return proxy.One2One, []proxy.Backend{backend.Relation()}, nil
    case tokenv1.TokenService_ServiceDesc.ServiceName:
      //return proxy.One2One, []proxy.Backend{backend.Token()}, nil
      return proxy.One2One, nil, status.Errorf(codes.Unimplemented, "Unknown method %s", methodName)
    }

    return proxy.One2One, nil, status.Errorf(codes.Unimplemented, "Unknown method %s", methodName)
  }

}
