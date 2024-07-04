package entity

import "errors"

type StorageProvider uint8

const (
  StorageProviderMinIO StorageProvider = iota
  StorageProviderAWSS3
  StorageProviderUnknown
)

func NewStorageProvider(val uint8) (StorageProvider, error) {
  provider := StorageProvider(val)
  if !provider.Valid() {
    return provider, errors.New("storage provider is not valid")
  }
  return provider, nil
}

func (s StorageProvider) Underlying() uint8 {
  return uint8(s)
}

func (s StorageProvider) Valid() bool {
  return s.Underlying() < StorageProviderUnknown.Underlying()
}

func (s StorageProvider) String() string {
  switch s {
  case StorageProviderMinIO:
    return "MinIO"
  case StorageProviderAWSS3:
    return "AWS S3"
  }
  return "Unknown"
}
