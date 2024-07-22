package external

import (
  "context"
  "github.com/arcorium/nexa/proto/gen/go/common"
  mailerv1 "github.com/arcorium/nexa/proto/gen/go/mailer/v1"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/sony/gobreaker"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/authentication/config"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/util"
  "sync"
)

func NewMailerClient(conn grpc.ClientConnInterface, conf *config.CircuitBreaker) external.IMailerClient {
  breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:         "nexa-mailer",
    MaxRequests:  conf.MaxRequest,
    Interval:     conf.ResetInterval,
    Timeout:      conf.OpenStateTimeout,
    IsSuccessful: util.IsGrpcConnectivityError,
  })

  return &mailerClient{
    client:    mailerv1.NewMailerServiceClient(conn),
    tagClient: mailerv1.NewTagServiceClient(conn),
    tracer:    util.GetTracer(),
    cb:        breaker,
    isReady:   false,
  }
}

type mailerClient struct {
  client    mailerv1.MailerServiceClient
  tagClient mailerv1.TagServiceClient
  tracer    trace.Tracer
  cb        *gobreaker.CircuitBreaker

  // Cache for tag id
  tagIdMutex             sync.RWMutex
  isReady                bool
  verificationEmailTagId string
  forgotPasswordTagId    string
}

func (m *mailerClient) send(ctx context.Context, request *mailerv1.SendMailRequest) error {
  _, err := m.client.Send(ctx, request)
  return err
}

func (m *mailerClient) getTagIds(ctx context.Context) error {
  m.tagIdMutex.RLock()

  if m.isReady {
    m.tagIdMutex.RUnlock()
    return nil
  }
  m.tagIdMutex.RUnlock()
  m.tagIdMutex.Lock()
  defer m.tagIdMutex.Unlock()

  result, err := m.cb.Execute(func() (interface{}, error) {
    return m.tagClient.Find(ctx, &common.PagedElementInput{
      Element: 0,
      Page:    0,
    })
  })
  if err != nil {
    return err
  }

  for _, tags := range result.(*mailerv1.FindTagResponse).Tags {
    switch tags.Name {
    case "Email Verification":
      m.verificationEmailTagId = tags.Id
    case "Reset Password":
      m.forgotPasswordTagId = tags.Id
    }
  }
  m.isReady = true
  return nil
}

func (m *mailerClient) SendEmailVerification(ctx context.Context, dto *dto.SendEmailVerificationDTO) error {
  ctx, span := m.tracer.Start(ctx, "MailerClient.SendEmailVerification")
  defer span.End()

  err := m.getTagIds(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  request := mailerv1.SendMailRequest{
    Recipients: []string{dto.Recipient.String()},
    Subject:    "Email Verification",
    BodyType:   mailerv1.BodyType_BODY_TYPE_PLAIN,
    Body:       "Verification Token: " + dto.Token,
    TagIds:     []string{m.verificationEmailTagId},
  }

  _, err = m.cb.Execute(func() (interface{}, error) {
    return nil, m.send(ctx, &request)
  })

  return err
}

func (m *mailerClient) SendForgotPassword(ctx context.Context, passwordDTO *dto.SendForgotPasswordDTO) error {
  ctx, span := m.tracer.Start(ctx, "MailerClient.SendForgotPassword")
  defer span.End()

  err := m.getTagIds(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  request := mailerv1.SendMailRequest{
    Recipients: []string{passwordDTO.Recipient.String()},
    Subject:    "Forgot Password",
    BodyType:   mailerv1.BodyType_BODY_TYPE_PLAIN,
    Body:       "Forgot Password Token: " + passwordDTO.Token,
    TagIds:     []string{m.forgotPasswordTagId},
  }

  _, err = m.cb.Execute(func() (interface{}, error) {
    return nil, m.send(ctx, &request)
  })

  return err
}
