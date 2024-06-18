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
  "nexa/services/mailer/internal/domain/repository"
  "nexa/services/mailer/util"
)

func NewSMTP(mailRepo repository.IMail, config *SMTPConfig) external.IMail {
  mail := &SMTPMailer{
    mailRepo: mailRepo,
    config:   *config,
    tracer:   util.GetTracer(),
    dialer:   gomail.NewDialer(config.Host, int(config.Port), config.Username, config.Password),
  }

  return mail
}

type SMTPConfig struct {
  Host      string
  Port      uint16
  Username  string
  Password  string
  TLSConfig *tls.Config
}

type SMTPMailer struct {
  mailRepo repository.IMail
  config   SMTPConfig
  dialer   *gomail.Dialer
  tracer   trace.Tracer
}

func (s *SMTPMailer) Send(ctx context.Context, mail *domain.Mail, attachments []dto.FileAttached) error {
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

  // Sending asynchronously
  go func() {
    err := s.dialer.DialAndSend(message)

    // Update status
    status := domain.StatusDelivered
    if err != nil {
      status = domain.StatusFailed
    }

    mail := domain.Mail{
      Id:     mail.Id,
      Status: status,
    }

    err = s.mailRepo.Patch(context.Background(), &mail)
  }()

  return nil
}

func CopyFunc(file dto.FileAttached) func(io.Writer) error {
  return func(writer io.Writer) error {
    _, err := writer.Write(file.Bytes)
    return err
  }
}
