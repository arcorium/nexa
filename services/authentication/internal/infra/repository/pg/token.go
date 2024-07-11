package pg

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/infra/repository/model"
  "nexa/services/authentication/util"
  "time"
)

func NewToken(db bun.IDB) repository.IToken {
  return &tokenRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type tokenRepository struct {
  db     bun.IDB
  tracer trace.Tracer
}

func (t *tokenRepository) Create(ctx context.Context, token *entity.Token) error {
  ctx, span := t.tracer.Start(ctx, "TokenRepository.Create")
  defer span.End()

  dbModel := model.FromTokenDomain(token, func(ent *entity.Token, token *model.Token) {
    token.CreatedAt = time.Now()
    token.ExpiredAt = ent.ExpiredAt
  })

  res, err := t.db.NewInsert().
    Model(&dbModel).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (t *tokenRepository) Delete(ctx context.Context, token string) error {
  ctx, span := t.tracer.Start(ctx, "TokenRepository.Delete")
  defer span.End()

  res, err := t.db.NewDelete().
    Model(types.Nil[model.Token]()).
    Where("token = ?", token).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (t *tokenRepository) DeleteByUserId(ctx context.Context, userId types.Id) error {
  ctx, span := t.tracer.Start(ctx, "TokenRepository.DeleteByUserId")
  defer span.End()

  res, err := t.db.NewDelete().
    Model(types.Nil[model.Token]()).
    Where("user_id = ?", userId.Underlying().String()).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (t *tokenRepository) Find(ctx context.Context, token string) (entity.Token, error) {
  ctx, span := t.tracer.Start(ctx, "TokenRepository.Find")
  defer span.End()

  var dbModel model.Token
  err := t.db.NewSelect().
    Model(&dbModel).
    Where("token = ?", token).
    OrderExpr("expired_at DESC").
    Scan(ctx)

  if err != nil {
    spanUtil.RecordError(err, span)
    return types.Null[entity.Token](), err
  }

  return dbModel.ToDomain()
}

func (t *tokenRepository) Get(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Token], error) {
  ctx, span := t.tracer.Start(ctx, "TokenRepository.Get")
  defer span.End()

  var dbModels []model.Token

  count, err := t.db.NewSelect().
    Model(&dbModels).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    ScanAndCount(ctx)

  result := repo.CheckPaginationResultWithSpan(dbModels, count, err, span)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return repo.NewPaginatedResult[entity.Token](nil, uint64(count)), result.Err
  }

  tokens, ierr := sharedUtil.CastSliceErrsP(result.Data, repo.ToDomainErr[*model.Token, entity.Token])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Token](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult(tokens, uint64(count)), nil
}
