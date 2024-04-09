package repository

import (
	"context"
	"github.com/uptrace/bun"
	"nexa/services/authentication/internal/domain/repository"
	"nexa/services/authentication/internal/infra/model"
	"nexa/services/authentication/shared/domain/entity"
	"nexa/shared/types"
	"nexa/shared/util"
	"nexa/shared/util/repo"
	"time"
)

func NewToken(db bun.IDB) repository.IToken {
	return &tokenRepository{db: db}
}

type tokenRepository struct {
	db bun.IDB
}

func (t *tokenRepository) Create(ctx context.Context, token *entity.Token) error {
	dbModel := model.FromTokenDomain(token, func(domain *entity.Token, token *model.Token) {
		token.CreatedAt = time.Now()
		token.ExpiredAt = domain.ExpiredAt
	})

	res, err := t.db.NewInsert().
		Model(&dbModel).
		Exec(ctx)

	return repo.CheckResult(res, err)
}

func (t *tokenRepository) Delete(ctx context.Context, token string) error {
	res, err := t.db.NewDelete().
		Model(util.Nil[model.Token]()).
		Where("token = ?", token).
		Exec(ctx)

	return repo.CheckResult(res, err)
}

func (t *tokenRepository) DeleteByUserId(ctx context.Context, userId types.Id) error {
	res, err := t.db.NewDelete().
		Model(util.Nil[model.Token]()).
		Where("user_id = ?", userId.Underlying().String()).
		Exec(ctx)

	return repo.CheckResult(res, err)
}

func (t *tokenRepository) Find(ctx context.Context, token string) (entity.Token, error) {
	var dbModel model.Token

	err := t.db.NewSelect().
		Model(&dbModel).
		Where("token = ?", token).
		Scan(ctx)

	return dbModel.ToDomain(), err
}

func (t *tokenRepository) FindAll(ctx context.Context, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Token], error) {
	var dbModels []model.Token

	count, err := t.db.NewSelect().
		Model(&dbModels).
		Limit(int(parameter.Limit)).
		Offset(int(parameter.Offset)).
		ScanAndCount(ctx)

	result := repo.CheckPaginationResult(dbModels, count, err)
	tokens := util.CastSlice(result.Data, func(from *model.Token) entity.Token {
		return from.ToDomain()
	})

	return repo.NewPaginatedResult(tokens, uint64(count)), result.Err
}
