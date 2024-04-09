package uow

import (
  "nexa/services/authentication/internal/domain/external"
  "sync"
)

const (
  userIndex       = 0
  tokenIndex      = 1
  tokenUsageIndex = 2
)

// ExternalClientManager repository storage, made members as private, so it couldn't be initialized outside this package
type ExternalClientManager struct {
  user external.IUserClient

  mutexes []sync.RWMutex
}

func (u *ExternalClientManager) User() external.IUserClient {
  u.mutexes[userIndex].RLock()
  return u.user
}

func (u *ExternalClientManager) UserDone() {
  u.mutexes[userIndex].RUnlock()
}

func (u *ExternalClientManager) SetUser(user external.IUserClient) {
  u.mutexes[userIndex].Lock()
  defer u.mutexes[userIndex].Unlock()

  u.user = user
}
