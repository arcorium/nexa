package entity

type StorageProvider uint8

const (
	StorageProviderMinIO StorageProvider = iota
	StorageProviderAWSS3
)

func (s StorageProvider) Underlying() uint8 {
	return uint8(s)
}

func (f StorageProvider) String() string {
	switch f {
	case StorageProviderMinIO:
		return "MinIO"
	case StorageProviderAWSS3:
		return "AWS S3"
	}
	return "Unknown"
}
