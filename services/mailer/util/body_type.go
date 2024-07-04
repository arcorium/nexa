package util

import (
  "errors"
  mailerv1 "nexa/proto/gen/go/mailer/v1"
  domain "nexa/services/mailer/internal/domain/entity"
)

func ToDomainBodyType(bodyType mailerv1.BodyType) (domain.MailBodyType, error) {
  switch bodyType {
  case mailerv1.BodyType_BODY_TYPE_PLAIN:
    return domain.BodyTypePlain, nil
  case mailerv1.BodyType_BODY_TYPE_HTML:
    return domain.BodyTypeHTML, nil
  }
  return domain.BodyTypeUnknown, errors.New("unknown body type")
}
