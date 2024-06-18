package dto

import (
  "nexa/shared/wrapper"
  "time"
)

type FileAttached struct {
  Filename string
  Bytes    []byte
}

type SendMailDTO struct {
  Subject   string `json:"subject" validate:"required"`
  Recipient string `json:"recipient" validate:"required,email"`
  Sender    wrapper.Nullable[string]
  BodyType  uint8    `json:"body_type" validate:"required"`
  Body      string   `json:"body"`
  TagIds    []string `json:"tag_ids"`

  Attachments []FileAttached
}

type UpdateMailDTO struct {
  Id            string   `json:"id" validate:"required,uuid4"`
  AddedTagIds   []string `json:"added_tag_ids" validate:"dive,uuid4"`
  RemovedTagIds []string `json:"removed_tag_ids" validate:"dive,uuid4"`
}

type MailResponseDTO struct {
  Id        string           `json:"id"`
  Subject   string           `json:"subject"`
  Recipient string           `json:"recipient"`
  Sender    string           `json:"sender"`
  Status    string           `json:"status"`
  SendAt    time.Time        `json:"send_at"` // TODO: Handle on mapper and service
  Tags      []TagResponseDTO `json:"tags"`
}
