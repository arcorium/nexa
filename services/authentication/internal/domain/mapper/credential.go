package mapper

import (
  "github.com/arcorium/nexa/shared/optional"
  "nexa/services/authentication/internal/domain/dto"
  "nexa/services/authentication/internal/domain/entity"
)

func ToCredentialResponseDTO(credential *entity.Credential, currentCredId optional.Object[string]) dto.CredentialResponseDTO {
  current := credential.Id.EqWithString(currentCredId.ValueOr(""))

  return dto.CredentialResponseDTO{
    Id:        credential.Id,
    Device:    credential.Device.Name,
    IsCurrent: current,
    // TODO: Add created and last refreshed
  }
}
