package external

import (
  "bytes"
  "context"
  "github.com/minio/minio-go/v7"
  "github.com/minio/minio-go/v7/pkg/credentials"
  "go.opentelemetry.io/otel/trace"
  "io"
  domain "nexa/services/file_storage/internal/domain/entity"
  "nexa/services/file_storage/internal/domain/external"
  "nexa/services/file_storage/util"
  "nexa/shared/types"
  spanUtil "nexa/shared/util/span"
  "time"
)

func NewMinIOStorage(useSSL bool, bucket string, config *MinIOClientConfig) (external.IStorage, error) {
  // Connect
  client, err := minio.New(config.Address, &minio.Options{
    Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
    Secure: useSSL,
  })

  if err != nil {
    return nil, err
  }

  obj := &MinIO{
    config: MinIOClientConfig{
      Address:         config.Address,
      AccessKeyID:     config.AccessKeyID,
      SecretAccessKey: config.SecretAccessKey,
    },
    Bucket: bucket,
    client: client,
    tracer: util.GetTracer(),
  }

  err = obj.setup()
  if err != nil {
    return nil, err
  }

  return obj, nil
}

type MinIOClientConfig struct {
  Address         string
  AccessKeyID     string
  SecretAccessKey string
}

type MinIO struct {
  Bucket string

  config MinIOClientConfig
  client *minio.Client
  tracer trace.Tracer
}

func (m *MinIO) Close(ctx context.Context) error {
  return nil
}

func (m *MinIO) GetProvider() domain.StorageProvider {
  return domain.StorageProviderMinIO
}

func (m *MinIO) setup() error {
  // Create bucket
  ctx := context.Background()
  exist, err := m.client.BucketExists(ctx, m.Bucket)
  if err != nil {
    return err
  }

  if !exist {
    // Create new one
    err = m.client.MakeBucket(ctx, m.Bucket, minio.MakeBucketOptions{})
    if err != nil {
      return err
    }

    // This policy allow the data to be read by anonymous
    policy := `
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Principal": {
            "AWS": [
              "*"
            ]
          },
          "Action": [
            "s3:GetObject"
          ],
          "Resource": [
            "arn:aws:s3:::` + m.Bucket + `/*"
          ]
        }
      ]
    }`
    err = m.client.SetBucketPolicy(ctx, m.Bucket, policy)
    if err != nil {
      return err
    }
  }
  return nil
}

func (m *MinIO) filepath(filename string) string {
  return filename
}

func (m *MinIO) Find(ctx context.Context, filename string) (domain.File, error) {
  ctx, span := m.tracer.Start(ctx, "MinIOStorage.Find")
  defer span.End()

  obj, err := m.client.GetObject(ctx, m.Bucket, filename, minio.GetObjectOptions{})
  if err != nil {
    spanUtil.RecordError(err, span)
    return domain.File{}, err
  }
  info, err := obj.Stat()

  // Read into bytes
  data, err := io.ReadAll(obj)
  if err != nil {
    spanUtil.RecordError(err, span)
    return domain.File{}, err
  }

  return domain.File{
    Name:  info.Key,
    Bytes: data,
    Size:  uint64(info.Size),
  }, nil
}

func (m *MinIO) Store(ctx context.Context, file *domain.File) (string, error) {
  ctx, span := m.tracer.Start(ctx, "MinIOStorage.Store")
  defer span.End()

  reader := bytes.NewReader(file.Bytes)
  opt := minio.PutObjectOptions{
    ContentType: util.GetMimeType(file.Name),
  }

  info, err := m.client.PutObject(ctx, m.Bucket, file.Name, reader, int64(file.Size), opt)
  if err != nil {
    spanUtil.RecordError(err, span)
    return "", err
  }

  return info.Key, nil
}

func (m *MinIO) Delete(ctx context.Context, filename string) error {
  ctx, span := m.tracer.Start(ctx, "MinIOStorage.Delete")
  defer span.End()

  err := m.client.RemoveObject(ctx, m.Bucket, filename, minio.RemoveObjectOptions{
    ForceDelete: true,
  })
  return err
}

func (m *MinIO) GetFullPath(ctx context.Context, filename string) (types.FilePath, error) {
  ctx, span := m.tracer.Start(ctx, "MinIOStorage.GetFullPath")
  defer span.End()

  url, err := m.client.PresignedGetObject(ctx, m.Bucket, m.filepath(filename), time.Hour*24, nil)
  if err != nil {
    spanUtil.RecordError(err, span)
    return "", err
  }

  return types.FilePathFromURL(url), nil
}

func (m *MinIO) Copy(ctx context.Context, src, dest string) (string, error) {
  ctx, span := m.tracer.Start(ctx, "MinIOStorage.Move")
  defer span.End()

  srcOpt := minio.CopySrcOptions{
    Bucket: m.Bucket,
    Object: src,
  }
  dstOpt := minio.CopyDestOptions{
    Bucket: m.Bucket,
    Object: dest,
  }

  // Copy the object into destination
  info, err := m.client.CopyObject(ctx, dstOpt, srcOpt)
  if err != nil {
    spanUtil.RecordError(err, span)
    return "", err
  }

  // Delete copied object (source)
  err = m.client.RemoveObject(ctx, m.Bucket, src, minio.RemoveObjectOptions{
    ForceDelete: true,
  })
  return info.Key, err
}
