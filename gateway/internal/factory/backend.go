package factory

import (
  "context"
  "github.com/siderolabs/grpc-proxy/proxy"
  "google.golang.org/grpc"
  "google.golang.org/grpc/metadata"
)

type Func func(ctx context.Context) (context.Context, *grpc.ClientConn, error)

type Backend struct {
  client *Client
}

func (b *Backend) prepare(ctx context.Context) context.Context {
  md, ok := metadata.FromIncomingContext(ctx)
  if ok {
    ctx = metadata.NewOutgoingContext(ctx, md.Copy())
  }
  return ctx
}

func (b *Backend) Authn() proxy.Backend {
  return &proxy.SingleBackend{GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
    newCtx := b.prepare(ctx)

    return newCtx, b.client.conn.Authn, nil
  }}
}

func (b *Backend) Authz() proxy.Backend {
  return &proxy.SingleBackend{GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
    newCtx := b.prepare(ctx)

    return newCtx, b.client.conn.Authz, nil
  }}
}

func (b *Backend) Comment() proxy.Backend {
  return &proxy.SingleBackend{GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
    newCtx := b.prepare(ctx)

    return newCtx, b.client.conn.Comment, nil
  }}
}

func (b *Backend) Feed() proxy.Backend {
  return &proxy.SingleBackend{GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
    newCtx := b.prepare(ctx)

    return newCtx, b.client.conn.Feed, nil
  }}
}

func (b *Backend) Storage() proxy.Backend {
  return &proxy.SingleBackend{GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
    newCtx := b.prepare(ctx)

    return newCtx, b.client.conn.Storage, nil
  }}
}

func (b *Backend) Mailer() proxy.Backend {
  return &proxy.SingleBackend{GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
    newCtx := b.prepare(ctx)

    return newCtx, b.client.conn.Mailer, nil
  }}
}

func (b *Backend) Post() proxy.Backend {
  return &proxy.SingleBackend{GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
    newCtx := b.prepare(ctx)

    return newCtx, b.client.conn.Post, nil
  }}
}

func (b *Backend) Reaction() proxy.Backend {
  return &proxy.SingleBackend{GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
    newCtx := b.prepare(ctx)

    return newCtx, b.client.conn.Reaction, nil
  }}
}

func (b *Backend) Relation() proxy.Backend {
  return &proxy.SingleBackend{GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
    newCtx := b.prepare(ctx)

    return newCtx, b.client.conn.Relation, nil
  }}
}

//func (b *Backend) Token() proxy.Backend {
//  return &proxy.SingleBackend{GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
//    newCtx := b.prepare(ctx)
//
//    return newCtx, b.client.conn.Token, nil
//  }}
//}

func NewBackend(client *Client) *Backend {
  return &Backend{
    client: client,
  }
}
