package mapper

import (
  "github.com/arcorium/nexa/proto/gen/go/common"
  postv1 "github.com/arcorium/nexa/proto/gen/go/post/v1"
  sharedDto "github.com/arcorium/nexa/shared/dto"
  sharedErr "github.com/arcorium/nexa/shared/errors"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "google.golang.org/protobuf/types/known/timestamppb"
  "nexa/services/post/internal/domain/dto"
  "nexa/services/post/internal/domain/entity"
)

func toProtoVisibility(visibility entity.Visibility) postv1.Visibility {
  switch visibility {
  case entity.VisibilityPublic:
    return postv1.Visibility_PUBLIC
  case entity.VisibilityFollower:
    return postv1.Visibility_FOLLOWER
  case entity.VisibilityOnlyMe:
    return postv1.Visibility_ONLY_ME
  default:
    return postv1.Visibility(entity.VisibilityUnknown)
  }
}

func ToEntityVisibility(visibility postv1.Visibility) (entity.Visibility, error) {
  switch visibility {
  case postv1.Visibility_PUBLIC:
    return entity.VisibilityPublic, nil
  case postv1.Visibility_FOLLOWER:
    return entity.VisibilityFollower, nil
  case postv1.Visibility_ONLY_ME:
    return entity.VisibilityOnlyMe, nil
  }
  return entity.VisibilityUnknown, sharedErr.ErrEnumOutOfBounds
}

func ToCreatePostDTO(request *postv1.CreatePostRequest) (dto.CreatePostDTO, error) {
  var fieldErrors []sharedErr.FieldError

  var sharedPostId *types.Id
  if request.SharedPostId != nil {
    id, err := types.IdFromString(*request.SharedPostId)
    if err != nil {
      fieldErrors = append(fieldErrors, sharedErr.NewFieldError("shared_post_id", err))
    } else {
      sharedPostId = &id
    }
  }

  visibility, err := ToEntityVisibility(request.Visibility)
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("visibility", err))
  }

  mediaIds, ierr := sharedUtil.CastSliceErrs(request.MediaIds, types.IdFromString)
  if !ierr.IsNil() {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("media_ids", ierr))
  }

  userIds, ierr := sharedUtil.CastSliceErrs(request.UserIds, types.IdFromString)
  if !ierr.IsNil() {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("user_ids", ierr))
  }

  if len(fieldErrors) > 0 {
    return dto.CreatePostDTO{}, sharedErr.GrpcFieldErrors2(fieldErrors...)
  }

  return dto.CreatePostDTO{
    SharedPostId: types.NewNullable(sharedPostId),
    Content:      types.NewNullable(request.Content),
    Visibility:   visibility,
    MediaIds:     mediaIds,
    UserIds:      userIds,
  }, nil
}

func ToEditPostDTO(request *postv1.EditPostRequest) (dto.EditPostDTO, error) {
  var fieldErrors []sharedErr.FieldError

  postId, err := types.IdFromString(request.PostId)
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("post_id", err))
  }

  mediaIds, ierr := sharedUtil.CastSliceErrs(request.MediaIds, types.IdFromString)
  if !ierr.IsNil() {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("media_ids", ierr))
  }

  userIds, ierr := sharedUtil.CastSliceErrs(request.UserIds, types.IdFromString)
  if !ierr.IsNil() {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("user_ids", ierr))
  }

  if len(fieldErrors) > 0 {
    return dto.EditPostDTO{}, sharedErr.GrpcFieldErrors2(fieldErrors...)
  }

  return dto.EditPostDTO{
    PostId:   postId,
    Content:  request.Content,
    MediaIds: mediaIds,
    UserIds:  userIds,
  }, nil
}

func toProtoTagUser(userDTO *dto.TaggedUserDTO) *postv1.UserTag {
  return &postv1.UserTag{
    UserId:   userDTO.UserId.String(),
    UserName: userDTO.Username,
  }
}

func ToProtoPost(respDTO *dto.PostResponseDTO) *postv1.Post {
  var parentPost *postv1.Post
  if respDTO.ParentPost != nil {
    parentPost = ToProtoPost(respDTO.ParentPost)
  }

  return &postv1.Post{
    Id:         respDTO.Id.String(),
    ParentPost: parentPost,
    CreatorId:  respDTO.CreatorId.String(),
    Content:    respDTO.Content,
    Visibility: toProtoVisibility(respDTO.Visibility),
    LastEdited: timestamppb.New(respDTO.LastEdited),
    CreatedAt:  timestamppb.New(respDTO.CreatedAt),
    Users:      sharedUtil.CastSliceP(respDTO.Tags, toProtoTagUser),
    MediaUrls:  respDTO.MediaUrls,
  }
}

func ToProtoEditedPost(responseDTO *dto.EditedPostResponseDTO) *postv1.EditedPost {
  posts := sharedUtil.CastSliceP(responseDTO.EditedPosts, func(from *dto.ChildPostResponseDTO) *postv1.EditedPost_Post {
    return &postv1.EditedPost_Post{
      Content:   from.Content,
      CreatedAt: timestamppb.New(from.CreatedAt),
      Users:     sharedUtil.CastSliceP(from.Tags, toProtoTagUser),
      MediaUrls: from.MediaUrls,
    }
  })

  return &postv1.EditedPost{
    PostId:    responseDTO.PostId.String(),
    CreatorId: responseDTO.CreatorId.String(),
    Posts:     posts,
  }
}

func ToProtoPagedOutput[T any](result *sharedDto.PagedElementResult[T]) *common.PagedElementOutput {
  return &common.PagedElementOutput{
    Element:       result.Element,
    Page:          result.Page,
    TotalElements: result.TotalElements,
    TotalPages:    result.TotalPages,
  }
}
