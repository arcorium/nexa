package uow

import "nexa/services/post/internal/domain/repository"

func NewStorage(post repository.IPost) PostStorage {
  return PostStorage{
    post: post,
  }
}

type PostStorage struct {
  post repository.IPost
}

func (m *PostStorage) Post() repository.IPost {
  return m.post
}
