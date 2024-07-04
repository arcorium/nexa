package external

import (
  "context"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc"
  mailerv1 "nexa/proto/gen/go/mailer/v1"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/external"
  "nexa/services/authentication/util"
)

func NewMailClient(conn grpc.ClientConnInterface) external.IMailClient {
  return &mailerClient{
    client: mailerv1.NewMailerServiceClient(conn),
    tracer: util.GetTracer(),
  }
}

type mailerClient struct {
  client mailerv1.MailerServiceClient

  tracer trace.Tracer
}

func (m *mailerClient) Send(ctx context.Context, dto *dto.SendVerificationEmailDTO) error {
  ctx, span := m.tracer.Start(ctx, "MailerClient.Send")
  defer span.End()

  request := mailerv1.SendMailRequest{
    Recipients: []string{dto.Recipient.String()},
    Subject:    "Email Verification",
    BodyType:   mailerv1.BodyType_BODY_TYPE_PLAIN,
    Body:       "Verification Token: " + dto.Token, // TODO: Change into url
    TagIds:     nil,                                // TODO: Use enum instead?
  }

  _, err := m.client.Send(ctx, &request)
  return err
}
