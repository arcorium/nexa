package mapper

import (
  "nexa/services/relation/internal/domain/dto"
  "nexa/services/relation/internal/domain/entity"
)

func ToBlockResponseDTO(block *entity.Block) dto.BlockResponseDTO {
  return dto.BlockResponseDTO{
    UserId:    block.BlockedId,
    CreatedAt: block.CreatedAt,
  }
}

func ToBlockCountResponseDTO(count *entity.BlockCount) dto.BlockCountResponseDTO {
  return dto.BlockCountResponseDTO{
    UserId: count.UserId,
    Total:  count.TotalBlocked,
  }
}
