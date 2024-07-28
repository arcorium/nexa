package mapper

import (
  commentv1 "github.com/arcorium/nexa/proto/gen/go/comment/v1"
  "github.com/arcorium/nexa/proto/gen/go/common"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "google.golang.org/protobuf/types/known/timestamppb"
  "nexa/services/comment/internal/domain/dto"
  "nexa/services/comment/internal/domain/entity"
)

func ToEntityItemType(val commentv1.Type) (entity.ItemType, error) {
  switch val {
  case commentv1.Type_POST_COMMENT:
    return entity.ItemPostComment, nil
  case commentv1.Type_COMMENT_REPLY:
    return entity.ItemCommentReply, nil
  }
  return entity.ItemUnknown, sharedErr.ErrEnumOutOfBounds
}

func ToCreateCommentDTO(request *commentv1.CreateCommentRequest) (dto.CreateCommentDTO, error) {
  var fieldErrors []sharedErr.FieldError

  postId, err := types.IdFromString(request.PostId)
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("post_id", err))
  }

  var parentId *types.Id
  if request.ParentId != nil {
    parentIds, err := types.IdFromString(*request.ParentId)
    if err != nil {
      fieldErrors = append(fieldErrors, sharedErr.NewFieldError("parent_id", err))
    }
    parentId = &parentIds
  }

  if len(fieldErrors) > 0 {
    return dto.CreateCommentDTO{}, sharedErr.GrpcFieldErrors2(fieldErrors...)
  }

  return dto.CreateCommentDTO{
    ParentId: types.NewNullable(parentId),
    PostId:   postId,
    Content:  request.Content,
  }, nil
}

func ToEditCommentDTO(request *commentv1.EditCommentRequest) (dto.EditCommentDTO, error) {
  commentId, err := types.IdFromString(request.CommentId)
  if err != nil {
    return dto.EditCommentDTO{}, sharedErr.NewFieldError("post_id", err).ToGrpcError()
  }

  return dto.EditCommentDTO{
    Id:      commentId,
    Content: request.Content,
  }, nil
}

func ToPagedElementDTO(input *common.PagedElementInput) sharedDto.PagedElementDTO {
  if input == nil {
    return sharedDto.PagedElementDTO{}
  }
  return sharedDto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
  }
}

func ToFindCommentByIdDTO(request *commentv1.FindCommentByIdRequest) (dto.FindCommentByIdDTO, error) {
  commentId, err := types.IdFromString(request.CommentId)
  if err != nil {
    return dto.FindCommentByIdDTO{}, sharedErr.NewFieldError("comment_id", err).ToGrpcError()
  }

  return dto.FindCommentByIdDTO{
    CommentId: commentId,
    ShowReply: request.ShowReply,
  }, nil
}

func ToGetPostsCommentsDTO(request *commentv1.GetPostCommentsRequest) (dto.GetPostsCommentsDTO, sharedDto.PagedElementDTO, error) {
  postId, err := types.IdFromString(request.PostId)
  if err != nil {
    return dto.GetPostsCommentsDTO{}, sharedDto.PagedElementDTO{}, sharedErr.NewFieldError("comment_id", err).ToGrpcError()
  }

  return dto.GetPostsCommentsDTO{
      PostId:    postId,
      ShowReply: request.ShowReply,
    },
    ToPagedElementDTO(request.Details),
    nil
}

func ToGetCommentRepliesDTO(request *commentv1.GetCommentRepliesRequest) (dto.GetCommentsRepliesDTO, sharedDto.PagedElementDTO, error) {
  postId, err := types.IdFromString(request.CommentId)
  if err != nil {
    return dto.GetCommentsRepliesDTO{}, sharedDto.PagedElementDTO{}, sharedErr.NewFieldError("comment_id", err).ToGrpcError()
  }

  return dto.GetCommentsRepliesDTO{
      CommentId: postId,
      ShowReply: request.ShowReply,
    },
    ToPagedElementDTO(request.Details),
    nil
}

func ToProtoComments(responseDTO *dto.CommentResponseDTO) *commentv1.Comment {
  var lastEdited *timestamppb.Timestamp
  if responseDTO.LastEdited.HasValue() {
    lastEdited = timestamppb.New(responseDTO.LastEdited.RawValue())
  }

  return &commentv1.Comment{
    Id:         responseDTO.Id.String(),
    PostId:     responseDTO.PostId.String(),
    UserId:     responseDTO.UserId.String(),
    Content:    responseDTO.Content,
    LastEdited: lastEdited,
    CreatedAt:  timestamppb.New(responseDTO.CreatedAt),
    Replies:    sharedUtil.CastSliceP(responseDTO.Replies, ToProtoComments),
  }
}

func ToProtoPagedElementOutput[T any](result *sharedDto.PagedElementResult[T]) *common.PagedElementOutput {
  return &common.PagedElementOutput{
    Element:       result.Element,
    Page:          result.Page,
    TotalElements: result.TotalElements,
    TotalPages:    result.TotalPages,
  }
}
