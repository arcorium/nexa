package handler

import (
  "context"
  "google.golang.org/protobuf/types/known/emptypb"
  "nexa/proto/gen/go/common"
  "nexa/proto/gen/go/mailer/v1"
  "nexa/services/mailer/internal/domain/service"
)

func NewMailHandler(mail service.IMail) MailHandler {
  return MailHandler{
    mailService: mail,
  }
}

type MailHandler struct {
  mailerv1.UnimplementedMailerServiceServer
  mailService service.IMail
}

func (m *MailHandler) Find(ctx context.Context, input *common.PagedElementInput) (*mailerv1.FindResponse, error) {
  //TODO implement me
  panic("implement me")
}

func (m *MailHandler) FindByIds(ctx context.Context, request *mailerv1.FindMailByIdsRequest) (*mailerv1.FindMailByIdsResponse, error) {
  //TODO implement me
  panic("implement me")
}

func (m *MailHandler) FindByTag(ctx context.Context, request *mailerv1.FindMailByTagRequest) (*mailerv1.FindMailByTagResponse, error) {
  //TODO implement me
  panic("implement me")
}

func (m *MailHandler) Send(ctx context.Context, request *mailerv1.SendMailRequest) (*mailerv1.SendMailResponse, error) {
  //TODO implement me
  panic("implement me")
}

func (m *MailHandler) SendToUser(ctx context.Context, request *mailerv1.SendMailToUserRequest) (*mailerv1.SendMailToUserResponse, error) {
  //TODO implement me
  panic("implement me")
}

func (m *MailHandler) Update(ctx context.Context, request *mailerv1.UpdateMailRequest) (*emptypb.Empty, error) {
  //TODO implement me
  panic("implement me")
}

func (m *MailHandler) Remove(ctx context.Context, request *mailerv1.RemoveMailRequest) (*emptypb.Empty, error) {
  //TODO implement me
  panic("implement me")
}

func (m *MailHandler) mustEmbedUnimplementedMailerServiceServer() {
  //TODO implement me
  panic("implement me")
}
