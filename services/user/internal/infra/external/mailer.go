package external

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  mailerv1 "nexa/proto/gen/go/mailer/v1"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/external"
)

func NewMailerClient(conn grpc.ClientConnInterface, tracer trace.Tracer) external.IMailerClient {
  return &mailerClient{
    client: mailerv1.NewMailerServiceClient(conn),
    tracer: tracer,
  }
}

type mailerClient struct {
  client mailerv1.MailerServiceClient

  tracer trace.Tracer
}

func (m *mailerClient) send(ctx context.Context, request *mailerv1.SendMailRequest) error {
  _, err := m.client.Send(ctx, request)
  return err
}

func (m *mailerClient) SendEmailVerification(ctx context.Context, dto *dto.SendEmailVerificationDTO) error {
  ctx, span := m.tracer.Start(ctx, "MailerClient.SendEmailVerification")
  defer span.End()

  request := mailerv1.SendMailRequest{
    Recipients: []string{dto.Recipient},
    Subject:    "Email Verification",
    BodyType:   mailerv1.BodyType_BODY_TYPE_PLAIN,
    Body:       "Verification Token: " + dto.Token,
    TagIds:     nil, // TODO: Use enum instead?
  }

  return m.send(ctx, &request)
}

func (m *mailerClient) SendForgotPassword(ctx context.Context, passwordDTO *dto.SendForgotPasswordDTO) error {
  ctx, span := m.tracer.Start(ctx, "MailerClient.SendForgotPassword")
  defer span.End()

  request := mailerv1.SendMailRequest{
    Recipients: []string{passwordDTO.Recipient},
    Subject:    "Email Verification",
    BodyType:   mailerv1.BodyType_BODY_TYPE_PLAIN,
    Body:       "Verification Token: " + passwordDTO.Token,
    TagIds:     nil, // TODO: Use enum instead?
  }

  return m.send(ctx, &request)
}
