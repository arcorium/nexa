package external

/*
func TestName(t *testing.T) {
  host := "sandbox.smtp.mailtrap.io"
  username := "b4b05a7f41be5d"
  password := "8612e5c9faf549"

  config := SMTPConfig{
    Host:     host,
    Port:     2525,
    Username: username,
    Password: password,
  }
  mailRepo := repository.NewMailMock(t)
  mailRepo.On("Patch", mock.Anything, mock.Anything).Return(nil)
  service, err := NewSMTP(mailRepo, &config)
  assert.NoError(t, err)

  mail := GenerateRandomMail()
  mail2 := GenerateRandomMail()

  err = service.Send(context.Background(), nil, mail, mail2)
  assert.NoError(t, err)
}

func GenerateRandomMail() domain.Mail {
  return domain.Mail{
    Id:        types.MustCreateId(),
    Subject:   sharedUtil.RandomString(20),
    Recipient: types.Email(gofakeit.Email()),
    Sender:    types.Email(gofakeit.Email()),
    BodyType:  domain.BodyTypePlain,
    Body:      "Hello",
    Status:    domain.StatusPending,
  }
}
*/
