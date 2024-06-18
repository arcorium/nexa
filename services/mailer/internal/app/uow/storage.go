package uow

import "nexa/services/mailer/internal/domain/repository"

func NewStorage(mail repository.IMail, tag repository.ITag) MailStorage {
  return MailStorage{
    mail: mail,
    tag:  tag,
  }
}

type MailStorage struct {
  mail repository.IMail
  tag  repository.ITag
}

func (m *MailStorage) Mail() repository.IMail {
  return m.mail
}

func (m *MailStorage) Tag() repository.ITag {
  return m.tag
}
