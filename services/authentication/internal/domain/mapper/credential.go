package mapper

import (
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
)

func ToCredentialResponseDTO(credential *entity.Credential) dto.CredentialResponseDTO {
  return dto.CredentialResponseDTO{
    Id:     credential.Id.Underlying().String(),
    Device: credential.Device.Name,
    // TODO: Add created and last refreshed
  }
}
