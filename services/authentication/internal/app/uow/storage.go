package uow

import (
  "nexa/services/authentication/internal/domain/repository"
)

func NewStorage(profile repository.IProfile, user repository.IUser) UserStorage {
  return UserStorage{
    profile: profile,
    user:    user,
  }
}

type UserStorage struct {
  profile repository.IProfile
  user    repository.IUser
}

func (m *UserStorage) Profile() repository.IProfile {
  return m.profile
}

func (m *UserStorage) User() repository.IUser {
  return m.user
}
