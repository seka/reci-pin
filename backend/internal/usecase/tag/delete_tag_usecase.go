package tag

import (
	"context"

	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

type DeleteTagUseCase interface {
	Execute(ctx context.Context, id int64) error
}

type deleteTagInteractor struct {
	tagRepo repository.TagRepository
}

func NewDeleteTagUseCase(tagRepo repository.TagRepository) DeleteTagUseCase {
	return &deleteTagInteractor{tagRepo: tagRepo}
}

func (uc *deleteTagInteractor) Execute(ctx context.Context, id int64) error {
	return uc.tagRepo.Delete(ctx, id)
}
