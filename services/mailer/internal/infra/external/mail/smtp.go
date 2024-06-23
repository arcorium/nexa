package mail

import (
  "context"
  "crypto/tls"
  "go.opentelemetry.io/otel/trace"
  "gopkg.in/gomail.v2"
  "io"
  "nexa/services/mailer/internal/domain/dto"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/services/mailer/internal/domain/external"
  "nexa/services/mailer/util"
)

func NewSMTP(config *SMTPConfig) (external.IMail, error) {
  dialer := gomail.NewDialer(config.Host, int(config.Port), config.Username, config.Password)
  sendCloser, err := dialer.Dial()
  if err != nil {
    return nil, err
  }

  return &SMTPMailer{
    config: *config,
    tracer: util.GetTracer(),
    dialer: dialer,
    sender: sendCloser,
  }, nil
}

type SMTPConfig struct {
  Host      string
  Port      uint16
  Username  string
  Password  string
  TLSConfig *tls.Config
}

type SMTPMailer struct {
  config SMTPConfig
  dialer *gomail.Dialer
  sender gomail.SendCloser // FIX: Make pool?
  tracer trace.Tracer
}

func (s *SMTPMailer) Send(ctx context.Context, mail *domain.Mail, attachments []dto.FileAttachment) error {
  ctx, span := s.tracer.Start(ctx, "SMTPMailer.Send")
  defer span.End()

  message := gomail.NewMessage(func(m *gomail.Message) {
    m.SetHeader("From", mail.Sender.Underlying())
    m.SetHeader("To", mail.Recipient.Underlying())
    m.SetHeader("Subject", mail.Subject)
    m.SetBody(mail.BodyType.String(), mail.Body)

    for _, attachment := range attachments {
      m.Attach(attachment.Filename, gomail.SetCopyFunc(CopyFunc(attachment)))
    }
  })

  return gomail.Send(s.sender, message)
}

func (s *SMTPMailer) Close(context.Context) error {
  return s.sender.Close()
}

func CopyFunc(file dto.FileAttachment) func(io.Writer) error {
  return func(writer io.Writer) error {
    _, err := writer.Write(file.Data)
    return err
  }
}
