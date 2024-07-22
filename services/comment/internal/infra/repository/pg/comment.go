package pg

import (
  "context"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "github.com/arcorium/nexa/shared/util/repo"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
  "github.com/uptrace/bun"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/comment/internal/domain/entity"
  "nexa/services/comment/internal/domain/repository"
  "nexa/services/comment/internal/infra/repository/model"
  "nexa/services/comment/util"
  "nexa/services/comment/util/errors"
)

func NewComment(db bun.IDB) repository.IComment {
  return &commentRepository{
    db:     db,
    tracer: util.GetTracer(),
  }
}

type commentRepository struct {
  db     bun.IDB
  tracer trace.Tracer
}

type idInput struct {
  Id string `bun:",type:uuid"`
}

func (c *commentRepository) GetReplies(ctx context.Context, showReply bool, commentId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Comment], error) {
  ctx, span := c.tracer.Start(ctx, "CommentRepository.GetReplies")
  defer span.End()

  var dbModels []model.Comment
  // Get children
  subQuery2 := c.db.NewSelect().
    Model(types.Nil[model.Comment]()).
    Join("JOIN result r ON comment.parent_id = r.id").
    OrderExpr("comment.created_at DESC")

  // Get parent
  subQuery := c.db.NewSelect().
    Model(types.Nil[model.Comment]()).
    Where("parent_id = ?", commentId.String()).
    OrderExpr("created_at DESC").
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset))

  if showReply {
    subQuery = subQuery.
      UnionAll(subQuery2)
  }

  // TODO: Get reply count for each reply
  err := c.db.NewSelect().
    WithRecursive("result", subQuery).
    Table("result").
    ColumnExpr("result.*").
    Scan(ctx, &dbModels)

  // Get total reply
  count, err := c.db.NewSelect().
    Model(types.Nil[model.Comment]()).
    Where("parent_id = ?", commentId.String()).
    Count(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  if result.Err != nil {
    spanUtil.RecordError(result.Err, span)
    return repo.NewPaginatedResult[entity.Comment](nil, uint64(count)), result.Err
  }

  comments, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Comment, entity.Comment])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Comment](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult(comments, uint64(count)), nil
}

func (c *commentRepository) GetReplyCounts(ctx context.Context, commentIds ...types.Id) ([]entity.Count, error) {
  ctx, span := c.tracer.Start(ctx, "CommentRepository.GetReplyCounts")
  defer span.End()

  var dbModels []model.Count
  input := sharedUtil.CastSlice(commentIds, func(id types.Id) idInput {
    return idInput{Id: id.String()}
  })
  err := c.db.NewSelect().
    With("result", c.db.NewValues(&input).WithOrder()).
    Model(types.Nil[model.Comment]()).
    ColumnExpr("result.id as ids").
    ColumnExpr("COUNT(parent_id) as total_comments").
    GroupExpr("result.id, result._order").
    OrderExpr("result._order").
    Join("RIGHT JOIN result ON comment.parent_id = result.id").
    Scan(ctx, &dbModels)

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  if len(dbModels) != len(commentIds) {
    // TODO: It should be internal error
    err = errors.ErrResultWithDifferentLength
    spanUtil.RecordError(err, span)
    return nil, err
  }

  counts := sharedUtil.CastSliceP(dbModels, repo.ToDomain[*model.Count, entity.Count])
  return counts, nil
}

func (c *commentRepository) GetPostCounts(ctx context.Context, postIds ...types.Id) ([]entity.Count, error) {
  ctx, span := c.tracer.Start(ctx, "CommentRepository.GetPostCounts")
  defer span.End()

  var dbModels []model.Count
  input := sharedUtil.CastSlice(postIds, func(id types.Id) idInput {
    return idInput{Id: id.String()}
  })
  err := c.db.NewSelect().
    With("result", c.db.NewValues(&input).WithOrder()).
    Model(types.Nil[model.Comment]()).
    ColumnExpr("result.id as ids").
    ColumnExpr("COUNT(post_id) as total_comments").
    GroupExpr("result.id, result._order").
    OrderExpr("result._order").
    Join("RIGHT JOIN result ON comment.post_id = result.id").
    Scan(ctx, &dbModels)

  if err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  if len(dbModels) != len(postIds) {
    // TODO: It should be internal error
    err = errors.ErrResultWithDifferentLength
    spanUtil.RecordError(err, span)
    return nil, err
  }

  counts := sharedUtil.CastSliceP(dbModels, repo.ToDomain[*model.Count, entity.Count])
  return counts, nil
}

func (c *commentRepository) FindByPostId(ctx context.Context, showReply bool, postId types.Id, parameter repo.QueryParameter) (repo.PaginatedResult[entity.Comment], error) {
  ctx, span := c.tracer.Start(ctx, "CommentRepository.FindByPostId")
  defer span.End()

  var dbModels []model.Comment
  // Get children
  subQuery2 := c.db.NewSelect().
    Model(types.Nil[model.Comment]()).
    //Relation("User").
    Join("JOIN result r ON comment.parent_id = r.id").
    Where("comment.post_id = ?", postId.String()).
    OrderExpr("comment.created_at DESC")

  // Get parent
  subQuery := c.db.NewSelect().
    Model(types.Nil[model.Comment]()).
    Where("parent_id IS NULL").
    Where("post_id = ?", postId.String()).
    OrderExpr("created_at DESC").
    Limit(int(parameter.Limit)).
    Offset(int(parameter.Offset))

  if showReply {
    subQuery = subQuery.
      UnionAll(subQuery2)
  }

  // TODO: Get count of reply
  err := c.db.NewSelect().
    WithRecursive("result", subQuery).
    Table("result").
    ColumnExpr("result.*").
    Scan(ctx, &dbModels)

  count, err := c.db.NewSelect().
    Model(types.Nil[model.Comment]()).
    Where("post_id = ? AND parent_id IS NULL", postId.String()). // Only take main comment
    Count(ctx)

  result := repo.CheckPaginationResult(dbModels, count, err)
  if result.Err != nil {
    spanUtil.RecordError(result.Err, span)
    return repo.NewPaginatedResult[entity.Comment](nil, uint64(count)), result.Err
  }

  comments, ierr := sharedUtil.CastSliceErrsP(dbModels, repo.ToDomainErr[*model.Comment, entity.Comment])
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return repo.NewPaginatedResult[entity.Comment](nil, uint64(count)), ierr
  }

  return repo.NewPaginatedResult(comments, uint64(count)), nil
}

func (c *commentRepository) Create(ctx context.Context, comment *entity.Comment) error {
  ctx, span := c.tracer.Start(ctx, "CommentRepository.Create")
  defer span.End()

  dbModel := model.FromCommentDomain(comment)
  res, err := c.db.NewInsert().
    Model(&dbModel).
    Returning("NULL").
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (c *commentRepository) UpdateContent(ctx context.Context, commentId types.Id, content string) error {
  ctx, span := c.tracer.Start(ctx, "CommentRepository.UpdateContent")
  defer span.End()

  res, err := c.db.NewUpdate().
    Model(types.Nil[model.Comment]()).
    Where("id = ?", commentId.String()).
    Set("content = ?", content).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (c *commentRepository) DeleteByIds(ctx context.Context, commentIds ...types.Id) error {
  ctx, span := c.tracer.Start(ctx, "CommentRepository.DeleteByIds")
  defer span.End()

  ids := sharedUtil.CastSlice(commentIds, sharedUtil.ToString[types.Id])
  // NOTE: Child comments deleted by on_cascade constraint
  res, err := c.db.NewDelete().
    Model(types.Nil[model.Comment]()).
    Where("id IN (?)", bun.In(ids)).
    Exec(ctx)

  return repo.CheckResultWithSpan(res, err, span)
}

func (c *commentRepository) DeleteUsers(ctx context.Context, userId types.Id, commentIds ...types.Id) ([]types.Id, error) {
  ctx, span := c.tracer.Start(ctx, "CommentRepository.DeleteUsers")
  defer span.End()

  // NOTE: Child comments deleted by on_cascade constraint
  query := c.db.NewDelete().
    Model(types.Nil[model.Comment]())

  if len(commentIds) == 0 {
    query = query.
      Where("user_id = ?", userId.String())
  } else {
    ids := sharedUtil.CastSlice(commentIds, sharedUtil.ToString[types.Id])
    query = query.
      Where("user_id = ? AND id IN (?)", userId.String(), bun.In(ids))
  }

  var output []idInput
  res, err := query.
    Returning("id").
    Exec(ctx, &output)

  if err = repo.CheckResult(res, err); err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  resp, ierr := sharedUtil.CastSliceErrsP(output, func(from *idInput) (types.Id, error) {
    return types.IdFromString(from.Id)
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }
  return resp, nil
}

func (c *commentRepository) DeleteByPostIds(ctx context.Context, postIds ...types.Id) ([]types.Id, error) {
  ctx, span := c.tracer.Start(ctx, "CommentRepository.DeleteByPostId")
  defer span.End()

  ids := sharedUtil.CastSlice(postIds, sharedUtil.ToString[types.Id])

  var output []idInput
  res, err := c.db.NewDelete().
    Model(types.Nil[model.Comment]()).
    Where("post_id IN (?)", bun.In(ids)).
    Returning("id").
    Exec(ctx, &output)

  if err = repo.CheckResult(res, err); err != nil {
    spanUtil.RecordError(err, span)
    return nil, err
  }

  resp, ierr := sharedUtil.CastSliceErrsP(output, func(from *idInput) (types.Id, error) {
    return types.IdFromString(from.Id)
  })
  if !ierr.IsNil() {
    spanUtil.RecordError(ierr, span)
    return nil, ierr
  }
  return resp, nil
}
