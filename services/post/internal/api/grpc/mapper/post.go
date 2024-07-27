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
  "time"
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

func ToPagedElementDTO(input *common.PagedElementInput) sharedDto.PagedElementDTO {
  if input == nil {
    return sharedDto.PagedElementDTO{
      Element: 0,
      Page:    0,
    }
  }
  return sharedDto.PagedElementDTO{
    Element: input.Element,
    Page:    input.Page,
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
  if !ierr.IsNil() && !ierr.IsEmptySlice() {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("media_ids", ierr))
  }

  userIds, ierr := sharedUtil.CastSliceErrs(request.UserIds, types.IdFromString)
  if !ierr.IsNil() && !ierr.IsEmptySlice() {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("user_ids", ierr))
  }

  if len(fieldErrors) > 0 {
    return dto.CreatePostDTO{}, sharedErr.GrpcFieldErrors2(fieldErrors...)
  }

  return dto.CreatePostDTO{
    SharedPostId:  types.NewNullable(sharedPostId),
    Content:       types.NewNullable(request.Content),
    Visibility:    visibility,
    MediaIds:      mediaIds,
    TaggedUserIds: userIds,
  }, nil
}

func ToEditPostDTO(request *postv1.EditPostRequest) (dto.EditPostDTO, error) {
  var fieldErrors []sharedErr.FieldError

  postId, err := types.IdFromString(request.PostId)
  if err != nil {
    fieldErrors = append(fieldErrors, sharedErr.NewFieldError("post_id", err))
  }

  var mediaIdsP *[]types.Id
  if !request.CloneLastMedia {
    mediaIds, ierr := sharedUtil.CastSliceErrs(request.MediaIds, types.IdFromString)
    if !ierr.IsNil() && !ierr.IsEmptySlice() {
      fieldErrors = append(fieldErrors, sharedErr.NewFieldError("media_ids", ierr))
    } else {
      mediaIdsP = &mediaIds
    }
  }

  var userIdsP *[]types.Id
  if !request.CloneLastTaggedUser {
    userIds, ierr := sharedUtil.CastSliceErrs(request.TaggedUserIds, types.IdFromString)
    if !ierr.IsNil() && !ierr.IsEmptySlice() {
      fieldErrors = append(fieldErrors, sharedErr.NewFieldError("tagged_user_ids", ierr))
    } else {

      userIdsP = &userIds
    }
  }

  if len(fieldErrors) > 0 {
    return dto.EditPostDTO{}, sharedErr.GrpcFieldErrors2(fieldErrors...)
  }

  return dto.EditPostDTO{
    PostId:   postId,
    Content:  request.Content,
    MediaIds: types.NewNullable(mediaIdsP),
    UserIds:  types.NewNullable(userIdsP),
  }, nil
}

func ToProtoPost(respDTO *dto.PostResponseDTO) *postv1.Post {
  var parentPost *postv1.Post
  if respDTO.ParentPost != nil {
    parentPost = ToProtoPost(respDTO.ParentPost)
  }
  var lastEdited *timestamppb.Timestamp
  if respDTO.LastEdited.Round(time.Second*2) != respDTO.CreatedAt.Round(time.Second*2) {
    lastEdited = timestamppb.New(respDTO.LastEdited)
  }

  return &postv1.Post{
    Id:         respDTO.Id.String(),
    ParentPost: parentPost,
    CreatorId:  respDTO.CreatorId.String(),
    Content:    respDTO.Content,
    Visibility: toProtoVisibility(respDTO.Visibility),
    LastEdited: lastEdited,
    CreatedAt:  timestamppb.New(respDTO.CreatedAt),
    TaggedUserIds: sharedUtil.CastSliceP(respDTO.Tags, func(tag *dto.TaggedUserDTO) string {
      return tag.UserId.String()
    }),
    MediaUrls: respDTO.MediaUrls,
  }
}

func ToProtoEditedPost(responseDTO *dto.EditedPostResponseDTO) *postv1.EditedPost {
  posts := sharedUtil.CastSliceP(responseDTO.EditedPosts, func(from *dto.ChildPostResponseDTO) *postv1.EditedPost_Post {
    return &postv1.EditedPost_Post{
      Content:       from.Content,
      CreatedAt:     timestamppb.New(from.CreatedAt),
      TaggedUserIds: sharedUtil.CastSliceP(from.Tags, func(tag *dto.TaggedUserDTO) string { return tag.UserId.String() }),
      MediaUrls:     from.MediaUrls,
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
