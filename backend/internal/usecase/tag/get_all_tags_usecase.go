package tag

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type GetAllTagsUseCase struct {
	tagRepo repository.TagRepository
}

func NewGetAllTagsUseCase(tagRepo repository.TagRepository) *GetAllTagsUseCase {
	return &GetAllTagsUseCase{tagRepo: tagRepo}
}

func (uc *GetAllTagsUseCase) Execute(ctx context.Context) ([]model.Tag, error) {
	return uc.tagRepo.GetAll(ctx)
}
