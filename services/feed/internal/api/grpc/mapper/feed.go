package mapper

import (
  feedv1 "github.com/arcorium/nexa/proto/gen/go/feed/v1"
  postv1 "github.com/arcorium/nexa/proto/gen/go/post/v1"
  "github.com/arcorium/nexa/shared/types"
  sharedUtil "github.com/arcorium/nexa/shared/util"
  "google.golang.org/protobuf/types/known/timestamppb"
  "nexa/services/feed/internal/domain/dto"
  "nexa/services/feed/internal/domain/entity"
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

func ToProtoPostWithCount(respDTO *dto.PostResponseDTO) *feedv1.PostWithCount {
  return &feedv1.PostWithCount{
    Post:     ToProtoPost(respDTO),
    Likes:    uint64(respDTO.TotalLikes),
    Dislikes: uint64(respDTO.TotalDislikes),
    Comments: uint64(respDTO.TotalComments),
  }
}

func ToProtoPost(respDTO *dto.PostResponseDTO) *postv1.Post {
  var parentPost *postv1.Post
  if respDTO.Parent != nil {
    parentPost = ToProtoPost(respDTO.Parent)
  }
  var lastEdited *timestamppb.Timestamp
  if respDTO.LastEditedAt.Round(time.Second*2) != respDTO.CreatedAt.Round(time.Second*2) {
    lastEdited = timestamppb.New(respDTO.LastEditedAt)
  }

  return &postv1.Post{
    Id:            respDTO.Id.String(),
    ParentPost:    parentPost,
    CreatorId:     respDTO.CreatorId.String(),
    Content:       respDTO.Content,
    Visibility:    toProtoVisibility(respDTO.Visibility),
    LastEdited:    lastEdited,
    CreatedAt:     timestamppb.New(respDTO.CreatedAt),
    TaggedUserIds: sharedUtil.CastSlice(respDTO.TaggedUserIds, sharedUtil.ToString[types.Id]),
    MediaUrls:     respDTO.MediaUrls,
  }
}
