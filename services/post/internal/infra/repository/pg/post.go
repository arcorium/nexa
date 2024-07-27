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
  "nexa/services/post/internal/domain/dto"
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

func (p *postRepository) getPost(models any) *bun.SelectQuery {
  return p.db.NewSelect().
    Model(models).
      Relation("Versions", func(query *bun.SelectQuery) *bun.SelectQuery {
        return query.
          Order("post_id", "created_at DESC").
          DistinctOn("post_id")
      }).
    Relation("Parent").
      Relation("Parent.Versions", func(query *bun.SelectQuery) *bun.SelectQuery {
        return query.
          Order("post_id", "created_at DESC").
          DistinctOn("post_id")
      }).
    Relation("Parent.Versions.Medias").
    Relation("Parent.Versions.UserTags").
    Relation("Versions.Medias").
    Relation("Versions.UserTags")
}

func (p *postRepository) Get(ctx context.Context, expectedVisibility entity.Visibility, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Post], error) {
  ctx, span := p.tracer.Start(ctx, "PostRepository.Get")
  defer span.End()

  var models []model.Post
  count, err := p.getPost(&models).
    Where("post.visibility <= ?", expectedVisibility.Underlying()).
    Limit(int(parameter.Limit)).
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
    Where("bookmark_post.user_id = ?", userId.String()).
    Relation("Post").
      Relation("Post.Versions", func(query *bun.SelectQuery) *bun.SelectQuery {
        return query.
          Order("post_id", "created_at DESC").
          DistinctOn("post_id")
      }).
    Relation("Post.Versions.Medias").
    Relation("Post.Versions.UserTags").
    Relation("Post.Parent").
      Relation("Post.Parent.Versions", func(query *bun.SelectQuery) *bun.SelectQuery {
        return query.
          Order("post_id", "created_at DESC").
          DistinctOn("post_id")
      }).
    Relation("Post.Parent.Versions.Medias").
    Relation("Post.Parent.Versions.UserTags").
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

func (p *postRepository) GetEdited(ctx context.Context, postId types.Id) (entity.Post, error) {
  ctx, span := p.tracer.Start(ctx, "PostRepository.GetEdited")
  defer span.End()

  var postModel model.Post
  err := p.db.NewSelect().
    Model(&postModel).
    Where("id = ?", postId.String()).
    Scan(ctx)

  if err != nil {
    spanUtil.RecordError(err, span)
    return entity.Post{}, err
  }

  var versionModels []model.PostVersion
  err = p.db.NewSelect().
    Model(&versionModels).
    Where("post_id = ?", postId.String()).
    Relation("Medias").
    Relation("UserTags").
    OrderExpr("created_at DESC").
    Scan(ctx)

  res := repo.CheckSliceResult(versionModels, err)
  if res.IsError() {
    spanUtil.RecordError(res.Err, span)
    return entity.Post{}, res.Err
  }

  postModel.Versions = versionModels

  post, err := postModel.ToDomain()
  if err != nil {
    spanUtil.RecordError(err, span)
    return entity.Post{}, err
  }

  return post, nil
}

func (p *postRepository) FindById(ctx context.Context, postIds ...types.Id) ([]entity.Post, error) {
  ctx, span := p.tracer.Start(ctx, "PostRepository.FindById")
  defer span.End()

  ids := sharedUtil.CastSlice(postIds, sharedUtil.ToString[types.Id])
  var models []model.Post
  err := p.getPost(&models).
    Where("post.id in (?)", bun.In(ids)).
    Scan(ctx)

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

func (p *postRepository) FindByUserId(ctx context.Context, userId types.Id, maxVisibility entity.Visibility, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Post], error) {
  ctx, span := p.tracer.Start(ctx, "PostRepository.FindByUserId")
  defer span.End()

  var models []model.Post
  count, err := p.getPost(&models).
    Where("post.creator_id = ? AND post.visibility <= ?", userId.String(), maxVisibility.Underlying()).
    Limit(int(parameter.Limit)).
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
    WherePK().
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

func (p *postRepository) insertPostVersions(ctx context.Context, db bun.IDB, versions ...model.PostVersion) error {
  // Insert post version
  res, err := db.NewInsert().
    Model(&versions).
    Returning("NULL").
    Exec(ctx)
  err = repo.CheckResult(res, err)
  if err != nil {
    return err
  }

  for _, version := range versions {
    // Insert user tags
    if len(version.UserTags) > 0 {
      res, err = db.NewInsert().
        Model(&version.UserTags).
        Returning("NULL").
        Exec(ctx)
      err = repo.CheckResult(res, err)
      if err != nil {
        return err
      }
    }

    // Insert medias
    if len(version.Medias) > 0 {
      res, err = db.NewInsert().
        Model(&version.Medias).
        Returning("NULL").
        Exec(ctx)
      err = repo.CheckResult(res, err)
      if err != nil {
        return err
      }
    }
  }
  return nil
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
      Model(&dbModel).
      Returning("NULL").
      Exec(ctx)
    err = repo.CheckResult(res, err)
    if err != nil {
      return err
    }

    return p.insertPostVersions(ctx, tx, dbModel.Versions...)
  })

  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  return nil
}

func (p *postRepository) UpdateVisibility(ctx context.Context, userId types.Id, postId types.Id, visibility entity.Visibility) error {
  ctx, span := p.tracer.Start(ctx, "PostRepository.UpdateVisibility")
  defer span.End()

  res, err := p.db.NewUpdate().
    Model(types.Nil[model.Post]()).
    Where("id = ? AND creator_id = ?", postId.String(), userId.String()).
    SetColumn("visibility", "?", visibility.Underlying()).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (p *postRepository) Edit(ctx context.Context, post *entity.Post, flag dto.EditPostFlag) error {
  ctx, span := p.tracer.Start(ctx, "PostRepository.Edit")
  defer span.End()

  // Get post existence
  exists, err := p.db.NewSelect().
    Model(types.Nil[model.Post]()).
    Where("id = ? AND creator_id = ?", post.Id.String(), post.CreatorId.String()).
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
  insertedVer := &dbModel.Versions[0]

  // Check if user tags or media ids is nil, if so then copy last versions
  if flag|dto.EditPostCopyNone != 0 {
    var models model.PostVersion
    query := p.db.NewSelect().
      Model(&models).
      Where("post_id = ?", post.Id.String()).
      Limit(1).
      OrderExpr("created_at DESC") // get newest

    if flag&dto.EditPostCopyTag != 0 {
      query = query.
        Relation("UserTags")
    }
    if flag&dto.EditPostCopyMedia != 0 {
      query = query.
        Relation("Medias")
    }

    err = query.Scan(ctx)
    if err != nil {
      spanUtil.RecordError(err, span)
      return err
    }

    // Set into dbModel
    for _, media := range models.Medias {
      insertedVer.Medias = append(insertedVer.Medias, model.Media{
        VersionId: insertedVer.Id,
        FileId:    media.FileId,
      })
    }
    for _, tag := range models.UserTags {
      insertedVer.UserTags = append(insertedVer.UserTags, model.UserTag{
        VersionId: insertedVer.Id,
        UserId:    tag.UserId,
      })
    }
  }

  err = p.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
    return p.insertPostVersions(ctx, tx, dbModel.Versions...)
  })
  if err != nil {
    spanUtil.RecordError(err, span)
    return err
  }

  return nil
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
