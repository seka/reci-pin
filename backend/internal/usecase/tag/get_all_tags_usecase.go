package tag

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type GetAllTagsUseCase interface {
	Execute(ctx context.Context) ([]model.Tag, error)
}

type getAllTagsInteractor struct {
	tagRepo repository.TagRepository
}

func NewGetAllTagsUseCase(tagRepo repository.TagRepository) GetAllTagsUseCase {
	return &getAllTagsInteractor{tagRepo: tagRepo}
}

func (uc *getAllTagsInteractor) Execute(ctx context.Context) ([]model.Tag, error) {
	return uc.tagRepo.GetAll(ctx)
}
