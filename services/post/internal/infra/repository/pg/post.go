package pg

import (
  "context"
  "database/sql"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/post/internal/domain/entity"
  "nexa/services/post/internal/domain/repository"
  "nexa/services/post/internal/infra/repository/model"
  "nexa/services/post/util"
  "time"
)

func NewPost(db bun.IDB) repository.IPost {
  return &postRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type postRepository struct {
  db     bun.IDB
  tracer trace.Tracer
}

func (p *postRepository) Get(ctx context.Context, child bool, expectedVisibility entity.Visibility, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Post], error) {
  ctx, span := p.tracer.Start(ctx, "PostRepository.Get")
  defer span.End()

  var models []model.Post
  query := p.db.NewSelect().
    Model(&models).
    Where("visibility <= ?", expectedVisibility.Underlying()).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    Relation("Versions").
    Relation("Versions.UserTags").
    Relation("Versions.Medias").
    OrderExpr("versions.created_at ASC")

  if !child {
    query.DistinctOn("versions.post_id")
  }

  count, err := query.Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    ScanAndCount(ctx)

  res := repo.CheckPaginationResult(models, count, err)
  if res.IsError() {
    spanUtil.RecordError(res.Err, span)
    return repo.NewPaginatedResult[entity.Post](nil, uint64(count)), res.Err
  }

  posts, ierr := sharedUtil.CastSliceErrsP(models, repo.ToDomainErr[*model.Post, entity.Post])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Post](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult[entity.Post](posts, uint64(count)), nil
}

func (p *postRepository) GetBookmarked(ctx context.Context, userId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Post], error) {
  ctx, span := p.tracer.Start(ctx, "PostRepository.GetBookmarked")
  defer span.End()

  var dbModels []model.BookmarkPost
  count, err := p.db.NewSelect().
    Model(&dbModels).
    Where("user_id = ?", userId.String()).
    Relation("Post").
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    ScanAndCount(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  if result.IsError() {
    spanUtil.RecordError(result.Err, span)
    return repo.NewPaginatedResult[entity.Post](nil, uint64(count)), result.Err
  }

  posts, ierr := sharedUtil.CastSliceErrsP(result.Data, func(from *model.BookmarkPost) (entity.Post, error) {
    return from.Post.ToDomain()
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Post](nil, uint64(count)), ierr
  }
  return repo.NewPaginatedResult(posts, uint64(count)), nil
}

func (p *postRepository) FindById(ctx context.Context, child bool, postIds ...types.Id) ([]entity.Post, error) {
  ctx, span := p.tracer.Start(ctx, "PostRepository.FindById")
  defer span.End()

  ids := sharedUtil.CastSlice(postIds, sharedUtil.ToString[types.Id])

  var models []model.Post
  query := p.db.NewSelect().
    Model(&models).
    Where("id in (?)", bun.In(ids)).
    Relation("Versions").
    Relation("Versions.UserTags").
    Relation("Versions.Medias").
    OrderExpr("versions.created_at ASC")

  if !child {
    query.DistinctOn("versions.post_id")
  }

  err := query.Scan(ctx)

  res := repo.CheckSliceResult(models, err)
  if res.IsError() {
    spanUtil.RecordError(res.Err, span)
    return nil, res.Err
  }

  posts, ierr := sharedUtil.CastSliceErrsP(models, repo.ToDomainErr[*model.Post, entity.Post])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return posts, nil
}

func (p *postRepository) FindByUserId(ctx context.Context, child bool, userId types.Id, maxVisibility entity.Visibility, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Post], error) {
  ctx, span := p.tracer.Start(ctx, "PostRepository.FindByUserId")
  defer span.End()

  var models []model.Post
  query := p.db.NewSelect().
    Model(&models).
    Where("user_id = ? AND visibility <= ?", userId, maxVisibility).
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    Relation("Versions").
    Relation("Versions.UserTags").
    Relation("Versions.Medias").
    OrderExpr("versions.created_at ASC")

  if !child {
    query.DistinctOn("versions.post_id")
  }

  count, err := query.Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset)).
    ScanAndCount(ctx)

  res := repo.CheckPaginationResult(models, count, err)
  if res.IsError() {
    spanUtil.RecordError(res.Err, span)
    return repo.NewPaginatedResult[entity.Post](nil, uint64(count)), res.Err
  }

  posts, ierr := sharedUtil.CastSliceErrsP(models, repo.ToDomainErr[*model.Post, entity.Post])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Post](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult[entity.Post](posts, uint64(count)), nil
}

func (p *postRepository) deleteBookmark(ctx context.Context, post *model.BookmarkPost) error {
  ctx, span := p.tracer.Start(ctx, "PostRepository.deleteBookmark")
  defer span.End()

  res, err := p.db.NewDelete().
    Model(post).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (p *postRepository) insertBookmark(ctx context.Context, post *model.BookmarkPost) error {
  ctx, span := p.tracer.Start(ctx, "PostRepository.insertBookmark")
  defer span.End()

  res, err := p.db.NewInsert().
    Model(post).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (p *postRepository) DelsertBookmark(ctx context.Context, userId, postId types.Id) error {
  ctx, span := p.tracer.Start(ctx, "PostRepository.DelsertBookmark")
  defer span.End()

  dbModel := model.BookmarkPost{
    UserId: userId.String(),
    PostId: postId.String(),
  }

  // Check existence
  exists, err := p.db.NewSelect().
    Model(&dbModel).
    WherePK().
    Exists(ctx)

  if err != nil {
    return err
  }

  if exists {
    return p.deleteBookmark(ctx, &dbModel)
  }
  return p.insertBookmark(ctx, &dbModel)
}

func (p *postRepository) Create(ctx context.Context, post *entity.Post) error {
  ctx, span := p.tracer.Start(ctx, "PostRepository.Create")
  defer span.End()

  dbModel := model.FromPostDomain(post, func(ent *entity.Post, post *model.Post) {
    post.CreatedAt = time.Now()
  })

  err := p.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
    // Insert post
    res, err := tx.NewInsert().
      Model(dbModel).
      Returning("NULL").
      Exec(ctx)
    err = repo.CheckResult(res, err)
    if err != nil {
      return err
    }

    // Insert post version
    res, err = tx.NewInsert().
      Model(&dbModel.Versions).
      Returning("NULL").
      Exec(ctx)
    err = repo.CheckResult(res, err)
    if err != nil {
      return err
    }

    for _, version := range dbModel.Versions {
      // Insert media
      res, err = tx.NewInsert().
        Model(&version.UserTags).
        Returning("NULL").
        Exec(ctx)
      err = repo.CheckResult(res, err)
      if err != nil {
        return err
      }

      // Insert user tags
      res, err = tx.NewInsert().
        Model(&version.Medias).
        Returning("NULL").
        Exec(ctx)
      err = repo.CheckResult(res, err)
      if err != nil {
        return err
      }
    }

    return nil
  })

  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  return nil
}

func (p *postRepository) UpdateVisibility(ctx context.Context, postId types.Id, visibility entity.Visibility) error {
  ctx, span := p.tracer.Start(ctx, "PostRepository.UpdateVisibility")
  defer span.End()

  dbModel := model.Post{
    Id: postId.String(),
    Visibility: sql.NullInt64{
      Int64: int64(visibility.Underlying()),
      Valid: true,
    },
  }
  res, err := p.db.NewUpdate().
    Model(&dbModel).
    WherePK().
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (p *postRepository) Edit(ctx context.Context, post *entity.Post) error {
  ctx, span := p.tracer.Start(ctx, "PostRepository.Edit")
  defer span.End()

  // Get post existence
  exists, err := p.db.NewSelect().
    Model(types.Nil[model.Post]()).
    Where("id = ?", post.Id.String()).
    Exists(ctx)
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  if !exists {
    err = sql.ErrNoRows
    spanUtil.RecordError(err, span)
    return err
  }

  // Create new post versions
  dbModel := model.FromPostDomain(post)

  err = p.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {

    // Insert post version
    res, err := tx.NewInsert().
      Model(&dbModel.Versions).
      Returning("NULL").
      Exec(ctx)
    err = repo.CheckResult(res, err)
    if err != nil {
      return err
    }

    for _, version := range dbModel.Versions {
      // Insert media
      res, err = tx.NewInsert().
        Model(&version.UserTags).
        Returning("NULL").
        Exec(ctx)
      err = repo.CheckResult(res, err)
      if err != nil {
        return err
      }

      // Insert user tags
      res, err = tx.NewInsert().
        Model(&version.Medias).
        Returning("NULL").
        Exec(ctx)
      err = repo.CheckResult(res, err)
      if err != nil {
        return err
      }
    }

    return nil
  })
  return err
}

func (p *postRepository) DeleteUsers(ctx context.Context, userId types.Id, postIds ...types.Id) ([]types.Id, error) {
  ctx, span := p.tracer.Start(ctx, "PostRepository.Delete")
  defer span.End()

  ids := sharedUtil.CastSlice(postIds, sharedUtil.ToString[types.Id])

  // All related data such as post_versioning, medias, and user tags will also be deleted because of the constraint
  var dbModels []model.Post
  query := p.db.NewDelete().
    Model(&dbModels)

  if len(ids) == 0 {
    query = query.
      Where("creator_id = ?", userId.String()).
      Returning("id")
  } else {
    query = query.
      Where("creator_id = ? AND id in (?)", userId.String(), bun.In(ids))
  }

  res, err := query.Exec(ctx)

  err = repo.CheckResult(res, err)
  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  deletedPostIds, ierr := sharedUtil.CastSliceErrsP(dbModels, func(from *model.Post) (types.Id, error) {
    return types.IdFromString(from.Id)
  })
  if !ierr.IsNil() && !ierr.IsEmptySlice() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }

  return deletedPostIds, nil
}

func (p *postRepository) Delete(ctx context.Context, postIds ...types.Id) error {
  ctx, span := p.tracer.Start(ctx, "PostRepository.Delete")
  defer span.End()

  ids := sharedUtil.CastSlice(postIds, sharedUtil.ToString[types.Id])

  // All related data such as post_versioning, medias, and user tags will also be deleted because of the constraint
  res, err := p.db.NewDelete().
    Model(types.Nil[model.Post]()).
    Where("id in (?)", bun.In(ids)).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}
