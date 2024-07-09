package external

import (
  "context"
  mailerv1 "github.com/arcorium/nexa/proto/gen/go/mailer/v1"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  "nexa/services/user/internal/domain/dto"
  "nexa/services/user/internal/domain/external"
  "nexa/services/user/util"
)

func NewMailerClient(conn grpc.ClientConnInterface) external.IMailerClient {
  return &mailerClient{
    client: mailerv1.NewMailerServiceClient(conn),
    tracer: util.GetTracer(),
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
    Recipients: []string{dto.Recipient.String()},
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
    Recipients: []string{passwordDTO.Recipient.String()},
    Subject:    "Email Verification",
    BodyType:   mailerv1.BodyType_BODY_TYPE_PLAIN,
    Body:       "Verification Token: " + passwordDTO.Token,
    TagIds:     nil, // TODO: Use enum instead?
  }

  return m.send(ctx, &request)
}
