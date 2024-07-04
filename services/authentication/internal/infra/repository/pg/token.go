package pg

import (
  "context"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/authentication/internal/domain/entity"
  "nexa/services/authentication/internal/domain/repository"
  "nexa/services/authentication/internal/infra/repository/model"
  "nexa/services/authentication/util"
  "nexa/shared/types"
  sharedUtil "nexa/shared/util"
  "nexa/shared/util/repo"
  spanUtil "nexa/shared/util/span"
  "time"
)

func NewToken(db bun.IDB) repository.IToken {
  return &tokenRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type tokenRepository struct {
  db bun.IDB

  tracer trace.Tracer
}

func (t *tokenRepository) Create(ctx context.Context, token *entity.Token) error {
  ctx, span := t.tracer.Start(ctx, "TokenRepository.Create")
  defer span.End()

  dbModel := model.FromTokenDomain(token, func(domain *entity.Token, token *model.Token) {
    token.CreatedAt = time.Now()
    token.ExpiredAt = domain.ExpiredAt
  })

  res, err := t.db.NewInsert().
    Model(&dbModel).
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
    Scan(ctx)

  if err != nil {
    spanUtil.RecordError(err, span)
    return types.Null[entity.Token](), err
  }

  return dbModel.ToDomain(), nil
}

func (t *tokenRepository) FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Token], error) {
  ctx, span := t.tracer.Start(ctx, "TokenRepository.FindAll")
  defer span.End()

  var dbModels []model.Token

  count, err := t.db.NewSelect().
    Model(&dbModels).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    ScanAndCount(ctx)

  result := repo.CheckPaginationResultWithSpan(dbModels, count, err, span)
  tokens := sharedUtil.CastSliceP(result.Data, func(from *model.Token) entity.Token {
    return from.ToDomain()
  })

  return repo.NewPaginatedResult(tokens, uint64(count)), result.Err
}
