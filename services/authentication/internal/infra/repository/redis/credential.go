package redis

import (
  "context"
  "database/sql"
  "errors"
  "fmt"
  "github.com/redis/go-redis/v9"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/infra/repository/model"
  "nexa/services/authentication/util"
  "nexa/shared/optional"
  spanUtil "nexa/shared/span"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/repo"
  "nexa/shared/variadic"
  "strconv"
  "strings"
  "time"
)

var defaultCredential = &CredentialConfig{
  CredentialNamespace: "cred",
  //CredentialIndicesNamespace: "creds",
  CredentialKeyCount: "creds-count",
  UserNamespace:      "user",
  MaxRetries:         5,
}

func NewCredential(client redis.UniversalClient, config *CredentialConfig) repository.ICredential {
  conf := optional.New(config).ValueOr(*defaultCredential)
  return &credentialRepository{
    config: &conf,
    client: client,
    tracer: util.GetTracer(),
  }
}

type CredentialConfig struct {
  CredentialNamespace string // All credentials
  //CredentialIndicesNamespace string // Credential indices
  CredentialKeyCount string // Store count of credentials
  UserNamespace      string // For lists of credentials each user

  MaxRetries int // maximum retry on Watch command
}

type credentialRepository struct {
  config *CredentialConfig
  client redis.UniversalClient
  //client *redis.Client

  tracer trace.Tracer
}

func (c *credentialRepository) credKey(id string) string {
  return fmt.Sprintf("%s:%s", c.config.CredentialNamespace, id)
}

func (c *credentialRepository) userKey(id string) string {
  return fmt.Sprintf("%s:%s", c.config.UserNamespace, id)
}

func (c *credentialRepository) scanModel(span trace.Span, id string, durCmd *redis.DurationCmd, cmd *redis.MapStringStringCmd) (model.Credential, error) {
  // Scan into model
  var models model.Credential
  if err := cmd.Scan(&models); err != nil {
    spanUtil.RecordError(err, span)
    return types.Null[model.Credential](), err
  }

  // Get expiration or TTL
  times := durCmd.Val()
  models.ExpiresAt = time.Unix(int64(times.Seconds()), 0)

  // Set id
  models.Id = id
  return models, nil
}

// watchRetry Run command Watch and retry it when it fails
func (c *credentialRepository) watchRetry(ctx context.Context, f func(tx *redis.Tx) error, keys ...string) error {
  var err error
  for i := 0; i < c.config.MaxRetries; i++ {
    err = c.client.Watch(ctx, f, keys...)
    if err == nil {
      break
    }
    // Failed due to the key is changed
    if err == redis.TxFailedErr {
      continue
    }
    return err
  }
  return nil
}

func (c *credentialRepository) Create(ctx context.Context, credential *entity.Credential) error {
  ctx, span := c.tracer.Start(ctx, "CredentialRepository.Create")
  defer span.End()

  models := model.FromCredentialModel(credential)
  credKey := c.credKey(models.Id)
  userKey := c.userKey(models.UserId)

  // Prevent updating created keys
  intCmd := c.client.Exists(ctx, credKey)
  if err := intCmd.Err(); err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  if intCmd.Val() != 0 {
    spanUtil.RecordError(repo.ErrAlreadyExists, span)
    return repo.ErrAlreadyExists
  }

  cmds, err := c.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
    // Create credentials
    pipe.HSet(ctx, credKey, &models)
    // Add cred to user namespace for indexing
    pipe.SAdd(ctx, userKey, credKey)
    // Count up credentials count
    pipe.Incr(ctx, c.config.CredentialKeyCount)
    // Set expiration time
    pipe.ExpireAt(ctx, credKey, models.ExpiresAt)
    return nil
  })

  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  for _, cmd := range cmds {
    if err := cmd.Err(); err != nil {
      spanUtil.RecordError(err, span)
      return err
    }
  }

  return nil
}

func (c *credentialRepository) Patch(ctx context.Context, credential *entity.Credential) error {
  ctx, span := c.tracer.Start(ctx, "CredentialRepository.Patch")
  defer span.End()

  models := model.FromCredentialModel(credential)
  credKey := c.credKey(models.Id)
  patched := models.OmitZero()

  tx := func(tx *redis.Tx) error {
    // Check if key already exist
    exists := tx.Exists(ctx, credKey)
    if err := exists.Err(); err != nil {
      return err
    }

    if exists.Val() == 0 {
      return sql.ErrNoRows
    }

    cmds, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
      // Update expiry time
      if !models.ExpiresAt.IsZero() {
        pipe.ExpireAt(ctx, credKey, models.ExpiresAt)
      }
      pipe.HSet(ctx, credKey, patched)
      return nil
    })

    if err != nil {
      return err
    }

    for _, cmd := range cmds {
      if err := cmd.Err(); err != nil {
        return err
      }
    }

    // Patch
    return nil
  }

  if err := c.watchRetry(ctx, tx, credKey); err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  return nil
}

func (c *credentialRepository) Delete(ctx context.Context, refreshTokenIds ...types.Id) error {
  ctx, span := c.tracer.Start(ctx, "CredentialRepository.Delete")
  defer span.End()

  credKeys := sharedUtil.CastSlice(refreshTokenIds, func(from types.Id) string {
    return c.credKey(from.Underlying().String())
  })

  var userKeys = make(map[string][]any) // redis-client does use []interface{}, instead of []string

  tx := func(tx *redis.Tx) error {
    // Find user ids related to refresh tokens
    for _, credKey := range credKeys {
      stringRes := tx.HGet(ctx, credKey, "user_id")
      if stringRes.Err() != nil || stringRes.Val() == "" {
        continue
      }
      userKeys[c.userKey(stringRes.Val())] = append(userKeys[c.userKey(stringRes.Val())], credKey)
    }

    // Delete credential
    intRes := tx.Del(ctx, credKeys...)
    if err := intRes.Err(); err != nil {
      spanUtil.RecordError(err, span)
      return err
    }

    if intRes.Val() == 0 {
      spanUtil.RecordError(sql.ErrNoRows, span)
      return sql.ErrNoRows
    }

    cmds, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
      // Delete from user
      for user, creds := range userKeys {
        pipe.SRem(ctx, user, creds...)
      }

      // Decrement count
      pipe.DecrBy(ctx, c.config.CredentialKeyCount, int64(len(refreshTokenIds)))

      return nil
    })

    if err != nil {
      return err
    }

    // Check command pipeline errors
    for _, cmd := range cmds {
      if err := cmd.Err(); err != nil {
        return err
      }
    }

    return nil
  }

  if err := c.watchRetry(ctx, tx, credKeys...); err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  return nil
}

func (c *credentialRepository) DeleteByUserId(ctx context.Context, userId types.Id, ids ...types.Id) error {
  ctx, span := c.tracer.Start(ctx, "CredentialRepository.DeleteByUserId")
  defer span.End()

  userKey := c.userKey(userId.Underlying().String())

  tx := func(tx *redis.Tx) error {
    vars := variadic.New(ids)

    credKeys := make([]string, 0)
    if !vars.HasValue() {
      // Get all credential ids
      var cursor uint64 = 0
      for {
        scanCmd := tx.SScan(ctx, userKey, cursor, "*", 0)
        if err := scanCmd.Err(); err != nil {
          return err
        }
        keys, newCursor := scanCmd.Val()
        credKeys = append(credKeys, keys...)
        if newCursor == 0 {
          break
        }
        cursor = newCursor
      }

      // Prevent delete user that has no credentials
      if len(credKeys) == 0 {
        return sql.ErrNoRows
      }
    } else {
      members := sharedUtil.CastSlice(ids, func(from types.Id) any {
        return c.credKey(from.Underlying().String())
      })
      boolRes := tx.SMIsMember(ctx, userKey, members...)
      if err := boolRes.Err(); err != nil {
        return err
      }

      if len(boolRes.Val()) != len(members) {
        return errors.New("malformed redis response")
      }

      // Filter
      for i := 0; i < len(boolRes.Val()); i++ {
        // Only append that is on users
        if boolRes.Val()[i] {
          credKeys = append(credKeys, members[i].(string))
        }
      }
    }

    cmds, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
      // Remove all creds
      pipe.Del(ctx, credKeys...)
      // Remove users sets
      //pipe.Del(ctx, userKey)

      return nil
    })

    if err != nil {
      return err
    }

    for _, cmd := range cmds {
      if err := cmd.Err(); err != nil {
        return err
      }
    }

    return nil
  }

  if err := c.watchRetry(ctx, tx, userKey); err != nil {
    spanUtil.RecordError(err, span)
    return err
  }
  return nil
}

func (c *credentialRepository) Find(ctx context.Context, refreshTokenId types.Id) (*entity.Credential, error) {
  ctx, span := c.tracer.Start(ctx, "CredentialRepository.Find")
  defer span.End()

  // Get data
  var expireAt *redis.DurationCmd
  var hResult *redis.MapStringStringCmd

  credKey := c.credKey(refreshTokenId.Underlying().String())
  tx := func(tx *redis.Tx) error {
    // Get hash map
    hResult = tx.HGetAll(ctx, credKey)
    if err := hResult.Err(); err != nil {
      return err
    }

    if len(hResult.Val()) == 0 {
      return sql.ErrNoRows
    }

    // Get expiration or TTL
    expireAt = tx.ExpireTime(ctx, credKey)
    if err := expireAt.Err(); err != nil {
      return err
    }

    return nil
  }

  err := c.watchRetry(ctx, tx, credKey)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  models, err := c.scanModel(span, refreshTokenId.Underlying().String(), expireAt, hResult)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  domainObj := models.ToDomain()
  return &domainObj, nil
}

func (c *credentialRepository) FindByUserId(ctx context.Context, userId types.Id) ([]entity.Credential, error) {
  ctx, span := c.tracer.Start(ctx, "CredentialRepository.FindByUserId")
  defer span.End()

  userKey := c.userKey(userId.Underlying().String())
  hResults := make(map[string]types.Pair[*redis.DurationCmd, *redis.MapStringStringCmd])

  tx := func(tx *redis.Tx) error {
    // Get credential ids
    credKeys := make([]string, 0)
    var cursor uint64 = 0
    for {
      scanCmd := tx.SScan(ctx, userKey, cursor, "*", 0)
      if err := scanCmd.Err(); err != nil {
        return err
      }
      keys, newCursor := scanCmd.Val()
      credKeys = append(credKeys, keys...)
      if newCursor == 0 {
        break
      }
      cursor = newCursor
    }

    if len(credKeys) == 0 {
      return sql.ErrNoRows
    }

    //Get the data for each credential
    _, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
      for _, credKey := range credKeys {
        // Get data
        res := pipe.HGetAll(ctx, credKey)
        // Get expiration or TTL
        expireAt := pipe.ExpireTime(ctx, credKey)

        hResults[credKey] = types.NewPair(expireAt, res)
      }
      return nil
    })

    return err
  }

  if err := c.watchRetry(ctx, tx, userKey); err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  var models []entity.Credential
  // Scan models
  for credKey, hResult := range hResults {
    // Duration
    if err := hResult.First.Err(); err != nil {
      spanUtil.RecordError(err, span)
      return nil, err
    }
    // Hash Result
    if err := hResult.Second.Err(); err != nil {
      spanUtil.RecordError(err, span)
      return nil, err
    }

    // Get credential id
    split := strings.Split(credKey, ":")
    if len(split) != 2 {
      err := errors.New("invalid credential key")
      spanUtil.RecordError(err, span)
      return nil, err
    }

    obj, err := c.scanModel(span, split[1], hResult.First, hResult.Second)
    if err != nil {
      spanUtil.RecordError(err, span)
      return nil, err
    }

    models = append(models, obj.ToDomain())
  }

  return models, nil
}

func (c *credentialRepository) FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Credential], error) {
  // NOTE: Offset is ignored, because redis doesn't support that

  ctx, span := c.tracer.Start(ctx, "CredentialRepository.FindAll")
  defer span.End()

  // Get all keys
  var cursor uint64 = 0
  prefix := fmt.Sprintf("%s:*", c.config.CredentialNamespace)
  credKeys := make([]string, 0)

  var credCount int64 = int64(parameter.Limit)
  for {
    scanCmd := c.client.Scan(ctx, cursor, prefix, credCount)
    if err := scanCmd.Err(); err != nil {
      spanUtil.RecordError(err, span)
      return repo.PaginatedResult[entity.Credential]{}, err
    }

    val, newCursor := scanCmd.Val()
    if len(val) > 0 {
      credKeys = append(credKeys, val...)
    }
    credCount = max(credCount-int64(len(val)), 0)
    if len(credKeys) >= int(parameter.Limit) && parameter.Limit != 0 {
      break
    }

    if newCursor == 0 {
      break
    }
    cursor = newCursor
  }

  results := make(map[string]types.Pair[*redis.DurationCmd, *redis.MapStringStringCmd])
  var countResult *redis.StringCmd

  tx := func(tx *redis.Tx) error {
    // Get all credentials data
    _, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
      for _, credKey := range credKeys {
        res := pipe.HGetAll(ctx, credKey)
        // Get expiration or TTL
        expireAt := pipe.ExpireTime(ctx, credKey)

        results[credKey] = types.NewPair(expireAt, res)
      }

      countResult = pipe.Get(ctx, c.config.CredentialKeyCount)

      return nil
    })

    return err
  }

  if err := c.watchRetry(ctx, tx, credKeys...); err != nil {
    spanUtil.RecordError(err, span)
    return repo.PaginatedResult[entity.Credential]{}, err
  }

  models := make([]entity.Credential, 0, len(credKeys))
  // Scan object
  for credKey, result := range results {
    // Duration
    if err := result.First.Err(); err != nil {
      spanUtil.RecordError(err, span)
      return repo.PaginatedResult[entity.Credential]{}, err
    }
    // Hash Result
    if err := result.Second.Err(); err != nil {
      spanUtil.RecordError(err, span)
      return repo.PaginatedResult[entity.Credential]{}, err
    }

    // Get credential id
    split := strings.Split(credKey, ":")
    if len(split) != 2 {
      err := errors.New("invalid credential key")
      spanUtil.RecordError(err, span)
      return repo.PaginatedResult[entity.Credential]{}, err
    }

    obj, err := c.scanModel(span, split[1], result.First, result.Second)
    if err != nil {
      spanUtil.RecordError(err, span)
      return repo.PaginatedResult[entity.Credential]{}, err
    }

    models = append(models, obj.ToDomain())
  }

  // Get count
  var count int64
  if err := countResult.Err(); err != nil {
    spanUtil.RecordError(err, span)
    return repo.PaginatedResult[entity.Credential]{}, err
  }
  countParsed, err := strconv.ParseInt(countResult.Val(), 10, 64)
  if err != nil {
    spanUtil.RecordError(err, span)
    return repo.PaginatedResult[entity.Credential]{}, err
  }
  count = countParsed

  return repo.NewPaginatedResult(models, uint64(count)), nil
}
