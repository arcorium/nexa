package mail

import (
  "context"
  "crypto/tls"
  "github.com/rs/zerolog/log"
  "go.opentelemetry.io/otel/trace"
  "gopkg.in/gomail.v2"
  "io"
  "nexa/services/mailer/internal/domain/dto"
  domain "nexa/services/mailer/internal/domain/entity"
  "nexa/services/mailer/internal/domain/external"
  "nexa/services/mailer/internal/domain/repository"
  "nexa/services/mailer/util"
  "sync"
  "sync/atomic"
)

func NewSMTP(mailRepo repository.IMail, config *SMTPConfig) (external.IMail, error) {
  dialer := gomail.NewDialer(config.Host, int(config.Port), config.Username, config.Password)
  sendCloser, err := dialer.Dial()
  if err != nil {
    return nil, err
  }

  return &SMTPMailer{
    mailRepo: mailRepo,
    config:   *config,
    tracer:   util.GetTracer(),
    dialer:   dialer,
    sender:   sendCloser,
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
  mailRepo repository.IMail
  config   SMTPConfig
  dialer   *gomail.Dialer
  sender   gomail.SendCloser // FIX: Make pool?
  tracer   trace.Tracer
}

func (s *SMTPMailer) Send(ctx context.Context, attachments []dto.FileAttachment, mails ...domain.Mail) error {
  ctx, span := s.tracer.Start(ctx, "SMTPMailer.Send")
  defer span.End()

  count := &atomic.Int64{}
  count.Store(int64(len(mails)))

  wg := sync.WaitGroup{}
  for _, mail := range mails {
    // Send each email async
    wg.Add(1)
    go func() {
      message := gomail.NewMessage(func(m *gomail.Message) {
        m.SetHeader("From", mail.Sender.Underlying())
        m.SetHeader("To", mail.Recipient.Underlying())
        m.SetHeader("Subject", mail.Subject)
        m.SetBody(mail.BodyType.String(), mail.Body)

        for _, attachment := range attachments {
          m.Attach(attachment.Filename, gomail.SetCopyFunc(CopyFunc(attachment)))
        }
      })

      log.Printf("Done!")
      wg.Done() // stop waiting until there
      err := gomail.Send(s.sender, message)

      // Update status
      status := domain.StatusDelivered
      if err != nil {
        status = domain.StatusFailed
      }

      mails := domain.Mail{
        Id:     mail.Id,
        Status: status,
      }

      // NOTE: use context.Background, because the parameter context could be already cancelled when this code will be executed
      err = s.mailRepo.Patch(context.Background(), &mails)
    }()
  }

  wg.Wait()

  return nil
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
