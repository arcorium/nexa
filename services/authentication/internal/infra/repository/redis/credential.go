package redis

import (
  "context"
  "fmt"
  "github.com/redis/go-redis/v9"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/infra/repository/model"
  spanUtil "nexa/shared/span"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/repo"
)

func NewCredential(client redis.UniversalClient) repository.ICredential {

}

type CredentialConfig struct {
  CredentialNamespace string // All credentials
  UserNamespace       string // For lists of credentials each user
}

type credentialRepository struct {
  config *CredentialConfig
  //client redis.UniversalClient
  client *redis.Client

  tracer trace.Tracer
}

func (c *credentialRepository) credKey(id string) string {
  return fmt.Sprintf("%s:%s", c.config.CredentialNamespace, id)
}

func (c *credentialRepository) userKey(id string) string {
  return fmt.Sprintf("%s:%s", c.config.UserNamespace, id)
}

func (c *credentialRepository) Create(ctx context.Context, credential *entity.Credential) error {
  ctx, span := c.tracer.Start(ctx, "CredentialRepository.Create")
  defer span.End()

  models := model.FromCredentialModel(credential)
  credId := c.credKey(models.Id)
  userId := c.userKey(models.UserId)

  _, err := c.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
    // Create credentials
    setRes := pipe.HSet(ctx, credId, &models)
    if err := setRes.Err(); err != nil {
      return err
    }

    // Add cred to user namespace for indexing
    intRes := pipe.SAdd(ctx, userId, credId)
    if err := intRes.Err(); err != nil {
      return err
    }

    // Set expiration time
    boolRes := pipe.ExpireAt(ctx, credId, models.ExpiresAt)
    return boolRes.Err()
  })

  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  return nil
}

func (c *credentialRepository) Patch(ctx context.Context, credential *entity.Credential) error {
  ctx, span := c.tracer.Start(ctx, "CredentialRepository.Patch")
  defer span.End()

  models := model.FromCredentialModel(credential)
  key := c.credKey(models.Id)
  patched := models.OmitZero()

  tx := func(tx *redis.Tx) error {
    exists := c.client.Exists(ctx, key)
    if err := exists.Err(); err != nil {
      return err
    }

    return tx.HSet(ctx, key, patched).Err()
  }

  err := c.client.Watch(ctx, tx, key)
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  return nil
}

func (c *credentialRepository) Delete(ctx context.Context, refreshTokenIds ...types.Id) error {
  ctx, span := c.tracer.Start(ctx, "CredentialRepository.Delete")
  defer span.End()

  ids := sharedUtil.CastSlice(refreshTokenIds, func(from types.Id) string {
    return from.Underlying().String()
  })

  intRes := c.client.Del(ctx, ids...)
  if err := intRes.Err(); err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  return nil
}

func (c *credentialRepository) DeleteByUserId(ctx context.Context, userId types.Id) error {
  ctx, span := c.tracer.Start(ctx, "CredentialRepository.DeleteByUserId")
  defer span.End()

  userKey := c.userKey(userId.Underlying().String())

  tx := func(tx *redis.Tx) error {
    // Get Credential Ids
    members := tx.SMembers(ctx, userKey)
    if err := members.Err(); err != nil {
      spanUtil.RecordError(err, span)
      return err
    }

    _, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
      credIds := members.Val()

      // Remove all creds
      delRes := pipe.Del(ctx, credIds...)
      if err := delRes.Err(); err != nil {
        return delRes.Err()
      }

      // Remove users sets
      delRes = pipe.Del(ctx, userKey)
      return delRes.Err()
    })

    return err
  }

  err := c.client.Watch(ctx, tx, userKey)
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  return nil
}

func (c *credentialRepository) Find(ctx context.Context, refreshTokenId types.Id) (entity.Credential, error) {
  ctx, span := c.tracer.Start(ctx, "CredentialRepository.Find")
  defer span.End()

  credKey := c.credKey(refreshTokenId.Underlying().String())
  res := c.client.HGetAll(ctx, credKey)
  if err := res.Err(); err != nil {
    spanUtil.RecordError(err, span)
    return entity.Credential{}, err
  }

  var models model.Credential
  err := res.Scan(&models)
  if err != nil {
    spanUtil.RecordError(err, span)
    return entity.Credential{}, err
  }

  return models.ToDomain(), nil
}

func (c *credentialRepository) FindByUserId(ctx context.Context, userId types.Id) ([]entity.Credential, error) {
  ctx, span := c.tracer.Start(ctx, "CredentialRepository.FindByUserId")
  defer span.End()

  userKey := c.userKey(userId.Underlying().String())

  var models []entity.Credential
  tx := func(tx *redis.Tx) error {
    // Get credential ids
    members := tx.SMembers(ctx, userKey)
    if err := members.Err(); err != nil {
      return err
    }

    // Get the data for each credential
    _, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
      for _, credKey := range members.Val() {
        res := c.client.HGetAll(ctx, credKey)
        if err := res.Err(); err != nil {
          return err
        }

        var obj model.Credential
        err := res.Scan(&obj)
        if err != nil {
          return err
        }
        models = append(models, obj.ToDomain())
      }
      return nil
    })

    return err
  }

  err := c.client.Watch(ctx, tx, userKey)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  return models, nil

}

func (c *credentialRepository) FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Credential], error) {
  ctx, span := c.tracer.Start(ctx, "CredentialRepository.FindAll")
  defer span.End()

  // Get all keys
  keys := c.client.Keys(ctx, fmt.Sprintf("%s:*", c.config.CredentialNamespace))
  if err := keys.Err(); err != nil {
    spanUtil.RecordError(err, span)
    return repo.PaginatedResult[entity.Credential]{}, nil
  }

  // Get all values
  credKeys := keys.Val()
  for _, credId := range credKeys {
    res := c.client.HGetAll(ctx, credId)
    if err := res.Err(); err != nil {
      return err
    }

    var obj model.Credential
    err := res.Scan(&obj)
    if err != nil {
      return err
    }
    models = append(models, obj.ToDomain())
  }

}
