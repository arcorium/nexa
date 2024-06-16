package mapper

import (
	"nexa/services/file_storage/internal/domain/dto"
	domain "nexa/services/file_storage/internal/domain/entity"
	"nexa/shared/types"
	"nexa/shared/wrapper"
)

func MapUpdateFileMetadataDTO(input *dto.UpdateFileMetadataDTO) (domain.FileMetadata, error) {
	id, err := types.IdFromString(input.Id)
	if err != nil {
		return domain.FileMetadata{}, err
	}

	md := domain.FileMetadata{
		Id: id,
	}

	wrapper.SetOnNonNull(&md.Name, input.Name)
	wrapper.SetOnNonNull(&md.IsPublic, input.IsPublic)

	return md, nil
}

func ToFileMetadataResponse(metadata *domain.FileMetadata) dto.FileMetadataResponseDTO {
	return dto.FileMetadataResponseDTO{
		Id:           metadata.Id.Underlying().String(),
		Name:         metadata.Name,
		Size:         metadata.Size,
		Path:         metadata.FullPath,
		CreatedAt:    metadata.CreatedAt,
		LastModified: metadata.LastModified,
	}
}
