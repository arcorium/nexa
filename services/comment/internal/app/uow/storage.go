package uow

import "nexa/services/comment/internal/domain/repository"

func NewStorage(comment repository.IComment) CommentStorage {
  return CommentStorage{
    comment: comment,
  }
}

type CommentStorage struct {
  comment repository.IComment
}

func (m *CommentStorage) Comment() repository.IComment {
  return m.comment
}
